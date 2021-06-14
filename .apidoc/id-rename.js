/*
 * Copyright (C) 2021 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of NethServer.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 */

//
// 1. Read file names from stdin one by line
// 2. Copy each file to a new destination.
// 3. The destination path is formed by the first argument passed to the
//    script and other path components read from the "$id" JSON Schema attribute.
//

const fs = require('fs');
const readline = require('readline');
const path = require('path');

const rl = readline.createInterface({
    input: process.stdin,
    crlfDelay: Infinity
});

const OUTDIR = process.argv[2]

rl.on('line', (filePath) => {
    if(fs.existsSync(filePath) != true) {
        console.log("Path not found: " + filePath)
        return
    }
    try {
        const data = JSON.parse(fs.readFileSync(filePath, 'utf8'))

        // Use the "$id" attribute contents to build the destination file name
        const idx = new URL([data["$id"]]);

        if(idx.protocol == 'urn:') {
            console.log("Skip " + filePath + ": it has URN protocol")
            return
        }

        var parts = idx.pathname.split('/').filter(p => p.length > 0)
        parts.reverse()

        var destName = parts[0]
        var dirName = ""
        if(parts.length > 1) {
            dirName = OUTDIR + "/" + parts[1]
        } else {
            dirName = OUTDIR
        }

        // Create the module directory
        if(!fs.existsSync(dirName)) {
            fs.mkdirSync(dirName)
        }

        fs.copyFile(filePath, dirName + "/" + destName, (err) => {
            if (err) throw err;
            console.log("Copied " + filePath + " to " + dirName + "/" + destName)
        })

    } catch (err) {
        console.error(err)
    }
});
