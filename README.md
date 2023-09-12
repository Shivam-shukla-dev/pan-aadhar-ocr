# pan-aadhar-ocr
Extract PAN and Aadhaar (UIDAI) card numbers from scanned images using OCR

Pre-requisites required for Mac

# HOW TO RUN ?

1. Pre-requisites required
- rew install tesseract
- export LIBRARY_PATH="/opt/homebrew/lib"
- export CPATH="/opt/homebrew/include"

2. go run cmd/main.go -a path_of_aadhar_card_image -p path_of_pan_card_image