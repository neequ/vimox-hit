#!/usr/bin/env sh
set -eu
repo=neequ/vimox-hit
github_mirror=${VIMOX_HIT_GITHUB_MIRROR:-}
case "$github_mirror" in
  ''|gh-proxy.com) ;;
  *) echo 'VIMOX_HIT_GITHUB_MIRROR must be gh-proxy.com or empty' >&2; exit 1 ;;
esac
github_url() {
  case "$github_mirror" in
    gh-proxy.com) printf 'https://gh-proxy.com/%s' "$1" ;;
    *) printf '%s' "$1" ;;
  esac
}
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$os" in darwin) os=darwin ;; linux) os=linux ;; *) echo "Unsupported OS: $os" >&2; exit 1 ;; esac
case "$arch" in x86_64|amd64) arch=amd64 ;; aarch64|arm64) arch=arm64 ;; *) echo "Unsupported architecture: $arch" >&2; exit 1 ;; esac
release=$(curl -fsSL -H 'Accept: application/vnd.github+json' "$(github_url "https://api.github.com/repos/$repo/releases/latest")")
tag=$(printf '%s' "$release" | sed -n 's/.*"tag_name":"\([^"]*\)".*/\1/p')
asset="vimox-hit_${tag}_${os}_${arch}.tar.gz"
base="https://github.com/$repo/releases/download/$tag"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
curl -fsSL "$(github_url "$base/$asset")" -o "$tmp/$asset"
curl -fsSL "$(github_url "$base/SHA256SUMS")" -o "$tmp/SHA256SUMS"
(cd "$tmp" && grep "  $asset\$" SHA256SUMS | sha256sum -c -)
tar -xzf "$tmp/$asset" -C "$tmp"
install_dir="${HOME}/.local/bin"
mkdir -p "$install_dir"
install -m 0755 "$tmp/vimox-hit" "$install_dir/vimox-hit"
"$install_dir/vimox-hit" version
printf 'Installed to %s. Ensure it is on PATH.\n' "$install_dir"
