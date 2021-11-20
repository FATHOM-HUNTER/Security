# Cloud storage implementation using Azure Services
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-3-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

This project demonstrates the usage of [Azure Storage SDK](https://github.com/Azure/azure-storage-blob-go) and [SQL Database on Azure](https://azure.microsoft.com/products/azure-sql/database/) along with JWT tokens for authentication and some other features.

Frontend can be found [here](https://github.com/Ovenoboyo/azure-cloud-storage-frontend)

## Features

- Username and Password based user authentication
- Authenticated API calls using JWT Tokens
- Upload files as Blobs to Azure Storage Containers
- Maintains per-user version history for each file uploaded
- List all files uploaded by a specific user after authentication
- Download / Delete files by version
- Good looking frontend (subjective)

## Implementation and screenshots

![login page](screenshots/login.png)
![register page](screenshots/register.png)
![dashboard page](screenshots/dashboard.png)

## How to build

Clone this repository along with submodules

```bash
git clone --recurse-submodules --remote-submodules
```

### Using Make

To run the app

```bash
make run
```

To simply generate a compiled binary

```bash
make build
```

### Manually

Install all dependencies for the website using [Yarn](https://yarnpkg.com/)

```bash
cd frontend
yarn install 
```

Compile and bundle the frontend

```bash
yarn build
```

Copy the generated bundle to static directory

```bash
mkdir ../static
cp dist/ ../static
```

then run

To generate a binary

```bash
go build
```

To run the app without generating a binary

```bash
go run main.go
```

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/ani20ket"><img src="https://avatars.githubusercontent.com/u/53856919?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Aniket Thorat</b></sub></a><br /><a href="#infra-ani20ket" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
    <td align="center"><a href="https://github.com/Ovenoboyo"><img src="https://avatars.githubusercontent.com/u/36789504?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Sahil Gupte</b></sub></a><br /><a href="https://github.com/Ovenoboyo/azure_cloud_storage/commits?author=Ovenoboyo" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/Mercyssh"><img src="https://avatars.githubusercontent.com/u/41297391?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Ohm</b></sub></a><br /><a href="#design-Mercyssh" title="Design">🎨</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
