# ShAp-Planner

ShAp-Planner is a **self-hosted app** for managing finances, tasks, and data within shared households.  
The app is fully open source, lightweight, and can run on small devices like Raspberry Pi or older computers.

**Backend:** Go  
**Frontend:** Android (Kotlin)  
**License:** Unlicense - complete freedom for everyone

---

## Summary
1. [Features](#-features)
2. [Configuration](#-configuration)
3. [Setup](#-setup)
4. [Contributing](#-contributing)
5. [License](#-license)

---

## ⚡ Features

- Multi-account support
- JWT-based login system
- Role-based access control (user/admin)
- Self-hosted, lightweight backend
- Configuration via environment variables
- Easy to extend with custom modules

---

## ⚙️ Configuration

### Environment Variables

| Variable       | Description                                           | Example          |
|----------------|-------------------------------------------------------|----------------|
| `JWT_SECRET`   | Secret used to sign JWT tokens                        | `superrandomsecret123` |

---

## 📝 Setup

1. Clone the repository:
```bash
git clone https://git.miaurizius.de/MiauRizius/shap-planner-backend.git
cd shap-planner-backend
````

2. Set environment variables:

```bash
export JWT_SECRET="your_super_random_secret"
```

3. Run the server:

```bash
go run main.go
```

---

## 🧩 Contributing

* Fork the repo
* Make changes
* Submit pull requests

We welcome bug fixes, new features, and documentation improvements.

---

## 📜 License

This project is licensed under **[The Unlicense](https://unlicense.org/)** – it is free for anyone to use, modify, and distribute for any purpose, without restrictions.