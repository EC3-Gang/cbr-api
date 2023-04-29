# cbr-api
An unofficial codebreaker.xyz API that is basically a wrapper that scrapes the website, parses the HTML and returns the data in a JSON format.

All API endpoints only accept GET requests.

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

## Roadmap
- [x] Add Swagger UI
- [x] Add an endpoint to get user info
- [] Add an endpoint to get recent submissions
- [] Add an endpoint to get recent submissions of a user
- [] Add an endpoint to get recent submissions of a user in a specific challenge
- [] Add an endpoint to get all problems
- [] Add an endpoint to get a specific problem (likely not possible)
- [] Add an endpoint to get all users (also likely not possible)