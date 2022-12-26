const pptr = require('puppeteer')
const net = require('net')

const APP_BASEURL = process.env.APP_BASEURL || 'http://nginx/'
const APP_HOST = process.env.APP_HOST || 'nginx'
const FLAG = process.env.FLAG || 'CTF{dummyflag}'

const sleep = d => new Promise(r => setTimeout(r, d));

const crawl = async (path = '') => {
  const url = APP_BASEURL + path

  console.log("crawl start:", url)

  // TODO: launch browsers on initialization and reuse like this
  // https://github.com/google/google-ctf/blob/85f182eac4472f58fb1ab3b708384331a1f4d3ca/2022/quals/web-postviewer-bot/bot.js#L96
  const browser = await pptr.launch({
    args: ["--no-sandbox"]
  })

  const context = await browser.createIncognitoBrowserContext()

  // create page
  const page = await context.newPage()

  await page.setCookie({
    name: "FLAG",
    value: FLAG,
    domain: APP_HOST,
    path: "/",
  });

  try {
    await page.goto(url, {
      timeout: 2000,
    });
    await sleep(2000);
  } catch (e) {
    console.log('error', e);
  }

  await page.close();
  await context.close();
  await browser.close();

  console.log("crawl end:", url)
}

const main = async () => {
  const server = net.createServer((socket) => {
    socket.on('data', (data) => {
      const path = data.toString().trim()
      crawl(path)
    })
  })
  server.listen(3000, '0.0.0.0')
}

main();
