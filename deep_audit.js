const fs = require('fs');
const path = require('path');

const swagger = JSON.parse(fs.readFileSync('docs/swagger.json', 'utf8'));
const frontendDir = path.resolve('frontend');

function getAllFiles(dirPath, arrayOfFiles) {
    const files = fs.readdirSync(dirPath);
    arrayOfFiles = arrayOfFiles || [];

    files.forEach(function (file) {
        const filePath = path.join(dirPath, file);
        try {
            if (fs.statSync(filePath).isDirectory()) {
                if (!['node_modules', '.nuxt', '.output', 'dist', '.git', 'playwright-report'].includes(file)) {
                    arrayOfFiles = getAllFiles(filePath, arrayOfFiles);
                }
            } else {
                if (['.vue', '.ts', '.js'].includes(path.extname(file))) {
                    arrayOfFiles.push(filePath);
                }
            }
        } catch (err) {
            // Skip inaccessible files/directories
            console.warn(`Warning: Could not access ${filePath}, skipping`);
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

        // 1. Literal match (cleaning placeholders)
        const baseSearch = epPath.replace(/\{[^}]+\}/g, '').replace(/:[^/]+/g, '').replace(/\/$/, '');

        // 2. Regex match for dynamic segments
        // Convert /api/v1/centers/{id}/holidays to regex pattern center.*holidays
        let regexPattern = epPath
            .replace('/api/v1', '')
            .replace(/\{[^}]+\}/g, '.*')
            .replace(/:[^/]+/g, '.*')
            .replace(/\//g, '.*');
        const re = new RegExp(regexPattern, 'i');

        const matchedFiles = fileContents.filter(f => {
            // Check literal base
            if (baseSearch !== '' && baseSearch !== '/api/v1' && f.content.includes(baseSearch)) return true;
            // Check regex
            if (re.test(f.content)) return true;
            return false;
        }).map(f => path.relative(frontendDir, f.path));

        results.push({
            path: epPath,
            method: method.toUpperCase(),
            summary: details.summary,
            integrated: matchedFiles.length > 0,
            matches: matchedFiles
        });
    }
}

fs.writeFileSync('deep_audit_results.json', JSON.stringify(results, null, 2), 'utf8');
console.log('Deep audit complete. Results saved to deep_audit_results.json');
