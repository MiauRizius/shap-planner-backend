# ShAp-Planner

ShAp-Planner is a **self-hosted app** for managing finances, tasks, and data within shared households.  
The app is fully open source, lightweight, and can run on small devices like Raspberry Pi or older computers.

**[Backend](https://git.miaurizius.de/MiauRizius/shap-planner-backend):** Go  
**[Frontend](https://git.miaurizius.de/MiauRizius/shap-planner-android):** Android (Kotlin)  
**[License](https://git.miaurizius.de/MiauRizius/shap-planner-backend/src/branch/main/LICENSE):** [CC0 1.0](https://creativecommons.org/publicdomain/zero/1.0/)

---

## Installation

### Docker  Compose (recommended)
1. Download docker-compose.yaml
````shell
$ curl -L https://git.miaurizius.de/MiauRizius/shap-planner-backend/raw/branch/main/docker-compose.yaml -o docker-compose.yaml
````
or create it yourself and enter the following content
````yaml
services:
  shap-planner:
    image: git.miaurizius.de/miaurizius/shap-planner-backend:latest
    container_name: shap-planner
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - SHAP_JWT_SECRET=SECURE_RANDOM_STRING # Must be at least 32 characters long
    volumes:
      - ./appdata:/appdata # To edit your configuration files
````

2. Start the container
````shell
$ docker compose up -d
````

3. Edit configuration as you like

### Build from source

1. Clone the repository:
```bash
git clone https://git.miaurizius.de/MiauRizius/shap-planner-backend.git
cd shap-planner-backend
````

2. Set environment variables:

```bash
export SHAP_JWT_SECRET="your_super_random_secret"
```

3. Run the server:

```bash
go run main.go
```

## Configuration
### Environment Variables

| Variable          | Description                                           | Example          |
|-------------------|-------------------------------------------------------|----------------|
| `SHAP_JWT_SECRET` | Secret used to sign JWT tokens                        | `superrandomsecret123` |

---

## License

This work is marked <a href="https://creativecommons.org/publicdomain/zero/1.0/">CC0 1.0</a>