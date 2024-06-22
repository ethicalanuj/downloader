# downloader

downloader is a command-line tool written in Go for downloading files from URLs specified in an input file.

## Installation
1. **Clone the repository:**
```
git clone https://github.com/ethicalanuj/downloader.git
cd downloader
```

2. **Build the executable:**
```
go build -o downloader downloader.go
```

## Usage
### Basic Usage
**For help**
```
./downloader -h
```
```
Usage: downloader [-h] [-l inputfile] [-o outputdir] [-v]

Options:
  -h              Display this help message.
  -l inputfile    Specify the input file containing URLs (default: js-urls.txt).
  -o outputdir    Specify the output directory (default: jsoutput-files).
  -v              Enable verbose output.
  ```

**Example :**
To download files listed in jsUrls.txt to the jsoutputFiles directory:
```
./downloader -l jsUrls.txt -o jsoutputFiles -v
```