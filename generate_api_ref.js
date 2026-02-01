const fs = require('fs');
const swagger = JSON.parse(fs.readFileSync('docs/swagger.json', 'utf8'));

let markdown = '# API 端點完整對照表\n\n**最後更新**：2026年1月31日\n**來源**：swagger.json\n\n---\n\n';

const categories = {};

for (const path in swagger.paths) {
    for (const method in swagger.paths[path]) {
        const details = swagger.paths[path][method];
        const tag = details.tags ? details.tags[0] : 'Uncategorized';
        if (!categories[tag]) categories[tag] = [];
        categories[tag].push({
            method: method.toUpperCase(),
            path: path,
            summary: details.summary || '',
            description: details.description || ''
        });
    }
}

const sortedTags = Object.keys(categories).sort();

sortedTags.forEach(tag => {
    markdown += `## ${tag}\n\n`;
    markdown += '| Method | Endpoint | Summary |\n';
    markdown += '|:---|:---|:---|\n';
    categories[tag].forEach(ep => {
        markdown += `| ${ep.method} | \`${ep.path}\` | ${ep.summary} |\n`;
    });
    markdown += '\n';
});

fs.writeFileSync('API_REFERENCE_NEW.md', markdown, 'utf8');
console.log('Successfully generated API_REFERENCE_NEW.md');
