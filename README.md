# Google Font Localizer

## Usage

1. Retrieve a google fonts css file
   - From google fonts website
     1. Go to https://fonts.google.com/ and select your font.
     2. Under the **Selected famities** sidenav, find the text aria with HTML code.
     3. Copy the url that begins with `https://fonts.googleapis.com/css2?family...` to your *clipboard*.
   - Or from another website
     1. Right click and inspect element.
     2. Look for this element`<link href="https://fonts.googleapis.com/css2?family=..." rel="stylesheet">` under the `<head>`
     3. Copy the url that begins with `https://fonts.googleapis.com/css2?family...` to your *clipboard*.
2. Save css file
   1. Create a directory for the fonts to live in.
      e.g.: `mkdir -p fonts/roboto`
   2. Create the `style.css` file there and paste the contents of your *clipboard* in there.
3. Run the binary `google-font-localizer` in the terminal in that directory.
