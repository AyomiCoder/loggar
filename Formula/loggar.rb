class Loggar < Formula
  desc "AI-powered log triage CLI"
  homepage "https://loggar.space"
  version "0.2.0"

  if Hardware::CPU.arm?
    url "https://github.com/AyomiCoder/loggar/releases/download/v0.2.0/loggar_darwin_arm64"
    sha256 "a000f900ec608e8c0832110329520af09b5097a58e4126e614e3e706f0e6ed18"
  else
    url "https://github.com/AyomiCoder/loggar/releases/download/v0.2.0/loggar_darwin_amd64"
    sha256 "e5710cc57c3c30e4ee39f907868167813c1669f577d8aee6ec76f21f03bcf600"
  end

  def install
    bin.install "loggar_darwin_#{Hardware::CPU.arch}" => "loggar"
  end

  test do
    system "#{bin}/loggar", "version"
  end
end
