class Loggar < Formula
  desc "AI-powered log triage CLI"
  homepage "https://loggar.dev"
  version "0.1.0"
  
  if Hardware::CPU.arm?
    url "https://github.com/AyomiCoder/loggar/releases/download/v0.1.0/loggar_darwin_arm64"
    sha256 "REPLACE_WITH_ACTUAL_SHA256_ARM64"
  else
    url "https://github.com/AyomiCoder/loggar/releases/download/v0.1.0/loggar_darwin_amd64"
    sha256 "REPLACE_WITH_ACTUAL_SHA256_AMD64"
  end

  def install
    bin.install "loggar_darwin_#{Hardware::CPU.arch}" => "loggar"
  end

  test do
    system "#{bin}/loggar", "version"
  end
end
