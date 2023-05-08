# cbr-api
An unofficial codebreaker.xyz API that is basically a wrapper that scrapes the website, parses the HTML and returns the data in a JSON format.

All API endpoints only accept GET requests.

## Credits
- [simonfalke (ryo simp)](https://github.com/simonfalke-01)
- [Canaris](https://github.com/DET171)


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
cd cbr-scraper
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
