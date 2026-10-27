// Everything put into `output` is available in the flow as ${output.<key>}.
console.log('[search-term.js] building the search query');

output.searchTerm = 'Wikipedia';
output.expectedResult = 'Wikipedia';

output.meta = {
  suite: 'regression',
  dataSource: 'runScript',
};

console.log('[search-term.js] query = ' + output.searchTerm);
