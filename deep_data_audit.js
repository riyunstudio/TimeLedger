const fs = require('fs');
const path = require('path');

const swagger = JSON.parse(fs.readFileSync('docs/swagger.json', 'utf8'));
const frontendDir = path.resolve('frontend');

function getAllFiles(dirPath, arrayOfFiles) {
    const files = fs.readdirSync(dirPath);
    arrayOfFiles = arrayOfFiles || [];

    files.forEach(function (file) {
        const fullPath = path.join(dirPath, file);
        if (!fs.existsSync(fullPath)) return;

        if (fs.statSync(fullPath).isDirectory()) {
            if (!['node_modules', '.nuxt', '.output', 'dist', '.git', '.output'].includes(file)) {
                arrayOfFiles = getAllFiles(fullPath, arrayOfFiles);
            }
        } else {
            if (['.vue', '.ts', '.js'].includes(path.extname(file))) {
                arrayOfFiles.push(fullPath);
            }
        }
    });

    return arrayOfFiles;
}

const allFiles = getAllFiles(frontendDir);
const fileContents = allFiles.map(f => ({ path: f, content: fs.readFileSync(f, 'utf8') }));

// Extract all properties from definitions in swagger
const apiProperties = {};
for (const defName in swagger.definitions) {
    const def = swagger.definitions[defName];
    if (def.properties) {
        apiProperties[defName] = Object.keys(def.properties);
    }
}

const results = [];

// Audit Endpoints and their usage of fields
for (const epPath in swagger.paths) {
    for (const method in swagger.paths[epPath]) {
        const details = swagger.paths[epPath][method];
        const res200 = details.responses['200'];
        let schemaRef = null;

        if (res200 && res200.schema) {
            if (res200.schema.$ref) {
                schemaRef = res200.schema.$ref.split('/').pop();
            } else if (res200.schema.items && res200.schema.items.$ref) {
                schemaRef = res200.schema.items.$ref.split('/').pop();
            }
        }

        const fields = schemaRef ? apiProperties[schemaRef] : [];

        // Find files that likely use this endpoint
        const baseSearch = epPath.replace(/\{[^}]+\}/g, '').replace(/:[^/]+/g, '').replace(/\/$/, '');
        const matchedFiles = fileContents.filter(f => baseSearch !== '' && f.content.includes(baseSearch));

        const fieldUsage = {};
        if (fields) {
            fields.forEach(field => {
                const usage = matchedFiles.filter(f => f.content.includes(field)).map(f => path.relative(frontendDir, f.path));
                fieldUsage[field] = usage;
            });
        }

        results.push({
            path: epPath,
            method: method.toUpperCase(),
            summary: details.summary,
            schema: schemaRef,
            expectedFields: fields,
            integrated: matchedFiles.length > 0,
            matches: matchedFiles.map(f => path.relative(frontendDir, f.path)),
            fieldUsage: fieldUsage
        });
    }
}

fs.writeFileSync('deep_data_audit_results.json', JSON.stringify(results, null, 2), 'utf8');
console.log('Deep data audit complete. Results saved to deep_data_audit_results.json');