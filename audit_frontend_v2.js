const fs = require('fs');
const path = require('path');

const swagger = JSON.parse(fs.readFileSync('docs/swagger.json', 'utf8'));
const frontendDir = path.resolve('frontend');

function getAllFiles(dirPath, arrayOfFiles) {
    const files = fs.readdirSync(dirPath);
    arrayOfFiles = arrayOfFiles || [];

    files.forEach(function (file) {
        if (fs.statSync(dirPath + "/" + file).isDirectory()) {
            if (file !== 'node_modules' && file !== '.nuxt' && file !== '.output' && file !== 'dist') {
                arrayOfFiles = getAllFiles(dirPath + "/" + file, arrayOfFiles);
            }
        } else {
            arrayOfFiles.push(path.join(dirPath, "/", file));
        }
    });

    return arrayOfFiles;
}

const allFiles = getAllFiles(frontendDir);
const fileContents = allFiles.map(f => ({ path: f, content: fs.readFileSync(f, 'utf8') }));

const results = [];

for (const epPath in swagger.paths) {
    for (const method in swagger.paths[epPath]) {
        const details = swagger.paths[epPath][method];

        let searchPath = epPath.replace(/\{[^}]+\}/g, '').replace(/:[^/]+/g, '').replace(/\/$/, '');
        if (searchPath === '' || searchPath === '/api/v1') {
            results.push({ path: epPath, method: method.toUpperCase(), integrated: 'SKIPPED', files: [] });
            continue;
        }

        let searchStr = searchPath;
        if (searchStr.startsWith('/api/v1')) {
            searchStr = searchStr.substring(7);
        }

        const foundInFiles = fileContents
            .filter(f => f.content.includes(searchStr))
            .map(f => path.relative(frontendDir, f.path));

        results.push({
            path: epPath,
            method: method.toUpperCase(),
            summary: details.summary,
            integrated: foundInFiles.length > 0,
            files: foundInFiles
        });
    }
}

fs.writeFileSync('integration_audit.json', JSON.stringify(results, null, 2), 'utf8');
console.log('Audit complete. Results saved to integration_audit.json');
