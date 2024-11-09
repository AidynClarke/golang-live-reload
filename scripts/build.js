const { spawn } = require('child_process')
const fs = require("fs")
const { join } = require('path')

const commandLineArgs = require('command-line-args')
const commandLineUsage = require('command-line-usage')

const optionDefinitions = [
  { name: 'inFile', alias: 'i', type: String, description: "Input file" },
  { name: 'outFile', alias: 'o', type: String, description: "Output file" },
  { name: 'help', alias: 'h', type: Boolean, description: "Shows this help menu" },
]

const options = commandLineArgs(optionDefinitions)

if (options.help || !options.inFile || !options.outFile) {
  const sections = [
    {
      header: 'Build a go executable',
      content: 'Usage: build -i <input file> -o <output file>'
    },
    {
      header: 'Options',
      optionList: optionDefinitions
    }
  ]
  const usage = commandLineUsage(sections)
  console.log(usage)
  process.exit(1)
}

/**
 * @param {string} inFile - Input file
 * @param {string} outFile - Output file
 * @returns {Promise<boolean>}
 */
function build(inFile, outFile) {
  return new Promise((res, rej) => {
    const child = spawn('go', ['build', '-o', outFile, inFile], { stdio: 'inherit' })
    child.on('close', (code) => {
      if (code !== 0) {
        return void rej(new Error(`Build failed with exit code ${code}`))
      }
      res(true)
    })

    child.on('error', console.error)
  })
}


const gopath = process.env.GOPATH
const gobin = process.env.GOBIN

let path

if (gobin) path = gobin
else if (gopath) path = join(gopath, 'bin')
else throw new Error('GOPATH not set')

build(options.inFile, `${path}/${options.outFile}`)
