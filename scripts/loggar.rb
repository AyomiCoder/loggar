class Loggar < Formula
  desc "AI-powered log triage CLI"
  homepage "https://loggar.dev"
  version "0.1.0"
  
  if Hardware::CPU.arm?
    url "https://github.com/AyomiCoder/loggar/releases/download/v0.1.0/loggar_darwin_arm64"
    sha256 "60f69f3e530f0e9e6a4bf5b19f239e50c054e0ee12dc76dc531aad95a72183c3"
  else
    url "https://github.com/AyomiCoder/loggar/releases/download/v0.1.0/loggar_darwin_amd64"
    sha256 "abf13064da336ff2a28b5f7a4fee34071f6e19e895020cc6df1212fb1d5f6d6f"
  end

  def install
    bin.install "loggar_darwin_#{Hardware::CPU.arch}" => "loggar"
  end

  test do
    system "#{bin}/loggar", "version"
  end
end
