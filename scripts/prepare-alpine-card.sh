#!/usr/bin/env bash
#
# Stage an Alpine SD card that boots straight into matrix.
#
# Alpine's Raspberry Pi release is a tarball extracted onto a FAT32 filesystem rather than
# an image written to a disk, so no dd and no root are needed. This builds the tree that
# belongs on the card in ./alpine-card; copy its contents to the card yourself. It never
# formats, partitions or writes to a block device.
#
# The result is self-contained. The board needs no keyboard, no screen and no network: it
# boots, starts matrix, and that is all.
#
#   ./scripts/prepare-alpine-card.sh
#   ./scripts/prepare-alpine-card.sh --components "core device clock"
#   ./scripts/prepare-alpine-card.sh --debug
#
# --components  what matrix starts. Defaults to core, device and every game.
# --hostname    the board's hostname. Defaults to matrix.
# --debug       log at debug level to matrix.log on the card. Off by default: nothing
#               rotates that file, so a console left running would eventually fill the
#               card.
#
set -euo pipefail

ALPINE_VERSION="${ALPINE_VERSION:-3.24.1}"
ALPINE_ARCH=armhf                                   # Alpine's name for ARMv6
ALPINE_TARBALL="alpine-rpi-${ALPINE_VERSION}-${ALPINE_ARCH}.tar.gz"
# The release series' own directory, v3.24 for 3.24.1, rather than latest-stable: that one
# holds only whatever is current, so a pinned version stops being downloadable the day
# Alpine cuts the next release.
ALPINE_BASE="https://dl-cdn.alpinelinux.org/alpine/v${ALPINE_VERSION%.*}/releases/${ALPINE_ARCH}"
# Checksum of the tarball named above, kept here rather than read from the .sha256 beside
# it on the mirror: that file travels with whatever the mirror serves, so it would confirm
# a replaced asset instead of rejecting it. A different ALPINE_VERSION needs its own.
ALPINE_SHA256="${ALPINE_SHA256:-1b32841873b4ff6b7a2f7247d65867545253bf9aa39a3c72be1c33eea9ab4ecd}"

# Everything is relative to the working directory rather than to the script, so this
# behaves the same run from a checkout, where `make package` writes matrix.tar.gz to the
# repository root, and from a directory holding nothing but a downloaded release.
WORK="${WORK:-$PWD/alpine-card}"
CACHE="${CACHE:-$PWD/.alpine-cache}"
MATRIX_TARBALL="${MATRIX_TARBALL:-$PWD/matrix.tar.gz}"

PI_HOSTNAME=matrix
DEBUG=0
# core and device, plus every game. gamepad and emulator are left out: gamepad is the web
# controller, unreachable without a network, and emulator is a development tool. Trim this
# with --components if memory gets tight.
COMPONENTS="core device demo zigzag yumyum clock draw blocks getout animate light rollupdice"
# Go's collector shares the single core with rendering, and the frame path allocates on
# every frame, so the default collection rate shows up as a periodic stutter. Trading
# memory for fewer collections is the right way round on a board that does nothing else.
GOGC=800

while [ $# -gt 0 ]; do
    case "$1" in
        --hostname)   PI_HOSTNAME="${2:?--hostname needs a name}"; shift 2 ;;
        --components) COMPONENTS="${2:?--components needs a list}"; shift 2 ;;
        --debug)      DEBUG=1; shift ;;
        -h|--help)  sed -n '3,21p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        *)          echo "unknown argument: $1" >&2; exit 2 ;;
    esac
done

if [ "$DEBUG" -eq 1 ]; then
    LOG_LEVEL=debug
else
    LOG_LEVEL=info
fi

say() { printf '\n\033[1m==> %s\033[0m\n' "$*"; }
die() { printf '\033[31merror: %s\033[0m\n' "$*" >&2; exit 1; }

for tool in curl tar sha256sum; do
    command -v "$tool" >/dev/null || die "$tool is required but not installed"
done
[ -f "$MATRIX_TARBALL" ] || die "$MATRIX_TARBALL not found.

Download matrix.tar.gz from the releases page into this directory, or build it from a
checkout:

    make package ARMV6_CC=<your arm musl gcc>"

# --- 1. Fetch Alpine, verified -----------------------------------------------------------

say "Fetching Alpine $ALPINE_VERSION for $ALPINE_ARCH"
mkdir -p "$CACHE"
if [ ! -f "$CACHE/$ALPINE_TARBALL" ]; then
    curl -fSL --progress-bar -o "$CACHE/$ALPINE_TARBALL" "$ALPINE_BASE/$ALPINE_TARBALL"
else
    echo "cached: $ALPINE_TARBALL"
fi

echo "$ALPINE_SHA256  $CACHE/$ALPINE_TARBALL" | sha256sum --check --quiet \
    || die "checksum mismatch on $ALPINE_TARBALL.
Delete $CACHE and retry. A version other than $ALPINE_VERSION needs ALPINE_SHA256 set to
that release's checksum as well."

# --- 2. Stage the card contents ----------------------------------------------------------

say "Staging card contents in $WORK"
rm -rf "$WORK"
mkdir -p "$WORK"
tar xzf "$CACHE/$ALPINE_TARBALL" -C "$WORK"

# --- 3. Overlay applied at boot ----------------------------------------------------------
#
# Alpine unpacks any *.apkovl.tar.gz on the boot media over / during boot. Only additive
# files go in here: replacing something the base system needs would break the boot.
#
# Everything the board needs is in this overlay, so setup-alpine never has to be run and
# nothing has to be committed with lbu. Each boot starts from the same known state.

say "Building $PI_HOSTNAME.apkovl.tar.gz"
OVL="$(mktemp -d)"
trap 'rm -rf "$OVL"' EXIT
# mktemp creates this 700. The overlay is unpacked over /, so that mode would land on the
# root directory itself and leave it unreadable to anything but root.
chmod 755 "$OVL"

mkdir -p "$OVL/etc"
echo "$PI_HOSTNAME" > "$OVL/etc/hostname"

# Without this marker the boot skips its whole default service set. The initramfs does:
#
#     if [ -f "$sysroot/etc/.default_boot_services" ] || ! [ -f "$ovl" ]; then
#             rc_add devfs sysinit; rc_add mdev sysinit; rc_add hwdrivers sysinit
#             rc_add modloop sysinit; rc_add modules boot; ...
#
# so supplying any apkovl at all, however small, silently disables devfs, mdev,
# hwdrivers and modloop unless the marker is present. Without modloop nothing under
# /lib/modules exists and no kernel module can load; without mdev, device nodes for
# hotplugged hardware are never created.
touch "$OVL/etc/.default_boot_services"

# USB-serial drivers, so the Arduino gets a /dev/ttyUSB* or /dev/ttyACM* whichever chip it
# carries: ch341 for the CH340 clones, ftdi_sio for FTDI, cdc-acm for boards with a native
# USB stack. USB HID is built into the kernel, so the controller needs nothing here.
#
# Nothing puts the USB port into gadget mode. It stays a host, because it has to drive the
# Arduino and the controllers; the OTG ethernet gadget would claim the port as a peripheral
# and then nothing downstream enumerates, hub included.
cat > "$OVL/etc/modules" <<'EOF'
ch341
ftdi_sio
cdc-acm
EOF

# Loopback. The components talk to core over gRPC on localhost, so 127.0.0.1 has to be
# reachable even with no real network. The initramfs's default service set does not
# include networking -- setup-alpine is what normally adds it -- so on an appliance that
# never runs the installer, lo stays down and every dial fails with "network is
# unreachable".
mkdir -p "$OVL/etc/network" "$OVL/etc/runlevels/boot"
cat > "$OVL/etc/network/interfaces" <<'EOF'
auto lo
iface lo inet loopback
EOF
ln -s /etc/init.d/networking "$OVL/etc/runlevels/boot/networking"

# Start matrix at boot. This is the only reason the card exists, so it is not optional.
mkdir -p "$OVL/etc/local.d" "$OVL/etc/runlevels/default"

if [ "$DEBUG" -eq 1 ]; then
    LOG_SETUP='BOOT=/media/mmcblk0p1
LOG="$BOOT/matrix.log"

# sync so the log survives having the power pulled, which is how this board gets turned
# off in practice.
mount -o remount,rw,sync "$BOOT" 2>/dev/null

# keep the previous boot'"'"'s log for comparison
[ -f "$LOG" ] && mv -f "$LOG" "$BOOT/matrix.log.prev" 2>/dev/null'
    LOG_REDIRECT='>> "$LOG" 2>&1'
else
    # No log file: nothing rotates it, and the boot media stays read-only, which is what
    # the card wants when the power is cut rather than shut down.
    LOG_SETUP=''
    LOG_REDIRECT='>/dev/null 2>&1'
fi

cat > "$OVL/etc/local.d/matrix.start" <<EOF
#!/bin/sh
# Generated by prepare-alpine-card.sh. Delete this file to stop matrix starting at boot.
$LOG_SETUP

# Belt and braces: the networking service should already have done this, and without it
# every gRPC dial to 127.0.0.1 fails with "network is unreachable".
ip link set lo up 2>/dev/null

# core reads themes/, fonts/ and images/ from the working directory and silently tolerates
# them missing, so starting anywhere else renders nothing.
cd /media/mmcblk0p1/matrix || exit 0

export GOGC=$GOGC
./matrix-linux-armv6 start --log-level $LOG_LEVEL $COMPONENTS $LOG_REDIRECT &

exit 0
EOF
chmod 755 "$OVL/etc/local.d/matrix.start"
echo "starts: matrix-linux-armv6 start --log-level $LOG_LEVEL $COMPONENTS"

# local.d scripts only run if the local service is enabled, and it is not part of the
# default set the initramfs adds.
ln -s /etc/init.d/local "$OVL/etc/runlevels/default/local"

tar czf "$WORK/$PI_HOSTNAME.apkovl.tar.gz" --owner=0 --group=0 -C "$OVL" .
echo "contents:"
tar tzf "$WORK/$PI_HOSTNAME.apkovl.tar.gz" | sed 's/^/  /'

# --- 4. Add matrix -----------------------------------------------------------------------

say "Adding matrix"
TMPM="$(mktemp -d)"
tar xzf "$MATRIX_TARBALL" -C "$TMPM"
mv "$TMPM/matrix-package" "$WORK/matrix"
rm -rf "$TMPM"
echo "matrix/ contains:"
ls "$WORK/matrix" | sed 's/^/  /'

# --- 5. Done -----------------------------------------------------------------------------

say "Staged $(du -sh "$WORK" | cut -f1) in $WORK"

cat <<EOF

Next
----
1. Format the card as a single FAT32 partition, marked bootable.

2. Copy the *contents* of this directory to the root of the card, not the directory
   itself:

       $WORK

3. Put the card in the board and power it on. matrix starts by itself.
EOF

if [ "$DEBUG" -eq 1 ]; then
    cat <<'EOF'

4. The log is matrix.log on the card, and the previous boot's is matrix.log.prev.
EOF
fi

echo
