desc "Builds ultralist-interactive for release"

Envs = [
  { goos: "darwin", arch: "386" },
  { goos: "darwin", arch: "amd64" },
  { goos: "linux", arch: "arm" },
  { goos: "linux", arch: "arm64" },
  { goos: "linux", arch: "386" },
  { goos: "linux", arch: "amd64" },
  { goos: "windows", arch: "386" },
  { goos: "windows", arch: "amd64" }
].freeze

Version = "1.7.0".freeze

task :build do
  `rm -rf dist/#{Version}`
  Envs.each do |env|
    ENV["GOOS"] = env[:goos]
    ENV["GOARCH"] = env[:arch]
    puts "Building #{env[:goos]} #{env[:arch]}"
    `GOOS=#{env[:goos]} GOARCH=#{env[:arch]} CGO_ENABLED=0 go build -v -o dist/#{Version}/uli`
    if env[:goos] == "windows"
      puts "Creating windows executable"
      `mv dist/#{Version}/uli dist/#{Version}/uli.exe`
      `cd dist/#{Version} && zip uli.zip uli.exe`
      puts "Removing windows executable"
      `rm -rf dist/#{Version}/uli.exe`
    else
      puts "Tarring #{env[:goos]} #{env[:arch]}"
      `cd dist/#{Version} && tar -czvf uli#{env[:goos]}_#{env[:arch]}.tar.gz uli`
      puts "Removing dist/#{Version}/uli"
      `rm -rf dist/#{Version}/uli`
    end
  end
end

task :lb do
  `rm -rf ./uli`
  le = { goos: "darwin", arch: "amd64" }
  ENV["GOOS"] = le[:goos]
  ENV["GOARCH"] = le[:arch]
  puts "Building #{le[:goos]} #{le[:arch]}"
  `GOOS=#{le[:goos]} GOARCH=#{le[:arch]} CGO_ENABLED=0 go build -v -o ./uli`
end

desc "Tests all the things"
task :test do
  system "go test ./..."
end

task default: :test
