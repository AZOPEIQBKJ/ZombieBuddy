// Raster exports of the original vector artwork; no third-party image input.
// node publishing/assets/render.cjs /path/to/sharp
const path = require('node:path');
const sharp = require(process.argv[2] || 'sharp');
(async () => {
  for (const [input, output, size] of [
    ['workshop-cover.svg', 'workshop-cover.png', 256],
    ['community-mark.svg', 'community-icon-256.png', 256],
    ['community-mark.svg', 'community-icon-128.png', 128],
  ]) {
    await sharp(path.join(__dirname, input)).resize(size, size).png()
      .toFile(path.join(__dirname, output));
    console.log(`${output}: ${size}x${size}`);
  }
})().catch(error => { console.error(error); process.exitCode = 1; });
