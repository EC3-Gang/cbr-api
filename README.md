# cbr-api
An unofficial codebreaker.xyz API that is basically a wrapper that scrapes the website, parses the HTML and returns the data in a JSON format.

All API endpoints only accept GET requests.

## Credits
- [simonfalke (ryo simp)](https://github.com/simonfalke-01) for the idea and go part
- [Canaris](https://github.com/DET171) did most the code though (no you did not my go code is longer than your js) (that's cap)

##### As of 384228c
```
───────────────────────────────────────────────────────────────────────────────
Language                 Files     Lines   Blanks  Comments     Code Complexity
───────────────────────────────────────────────────────────────────────────────
JavaScript                   9       729       54        42      633         28
JSON                         2        50        0         0       50          0
gitignore                    2        20        4         3       13          0
Go                           1       200       20        10      170         41
License                      1         8        4         0        4          0
Markdown                     1        37        8         0       29          0
YAML                         1         1        0         0        1          0
───────────────────────────────────────────────────────────────────────────────
Total                       17      1045       90        55      900         69
───────────────────────────────────────────────────────────────────────────────
```


## API Endpoints
### Problems
#### `/api/getProblem`
##### Parameters
- `problemId`: `string` (required)

#### `/api/getProblems`
##### Parameters
- No parameters

### Submissions
#### `/api/getSubmissions`
##### Parameters
- `problemId`: `string` (required)
- `ac`: `bool` (default: false)

### Users
#### `/api/getUser`
##### Parameters
- `name`: `string` (required)

### Contests
#### `/api/getContest`
##### Parameters
- `contestId`: `string` (required)

#### `/api/getContests`
##### Parameters
- No parameters



## Usage

### Installation
```bash
git clone https://github.com/EC3-Gang/cbr-api.git
cd cbr-api
yarn
```

### Development
```bash
yarn dev
```
Navigate to `http://localhost:3000/`.

### Production
```bash
yarn start
```
Navigate to `http://localhost:3000/`.

### Lint
```bash
yarn lint
```

## Endpoints
Visit the root url to find out more about the endpoints.
