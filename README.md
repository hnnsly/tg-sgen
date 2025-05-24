```
,---------.    .-_'''-.               .-'''-.    .-_'''-.        .-''-.   ,---.   .--.
\          \  '_( )_   \             / _     \  '_( )_   \     .'_ _   \  |    \  |  |
 `--.  ,---' |(_ o _)|  '           (`' )/`--' |(_ o _)|  '   / ( ` )   ' |  ,  \ |  |
    |   \    . (_,_)/___|          (_ o _).    . (_,_)/___|  . (_ o _)  | |  |\_ \|  |
    :_ _:    |  |  .-----.          (_,_). '.  |  |  .-----. |  (_,_)___| |  _( )_\  |
    (_I_)    '  \  '-   .'         .---.  \  : '  \  '-   .' '  \   .---. | (_ o _)  |
   (_(=)_)    \  `-'`   |          \    `-'  |  \  `-'`   |   \  `-'    / |  (_,_)\  |
    (_I_)      \        /           \       /    \        /    \       /  |  |    |  |
    '---'       `'-...-'             `-...-'      `'-...-'      `'-..-'   '--'    '--'
```

Semi-automated Telegram session creation tool built with Gogram.

## 🚨 Important Warnings

**⚠️ USE YOUR OWN App ID and App Hash!**

- Get them at [my.telegram.org](https://my.telegram.org/auth)
- DO NOT use default/shared credentials
- Using others' credentials violates Telegram ToS and may result in bans

**📋 Disclaimer:**

- Use at your own risk
- Author is not responsible for account bans or ToS violations
- Comply with Telegram's Terms of Service
- Do not use for spam or abuse

## 🛠 Quick Start

### Docker

```bash
# Prepare phones.txt with phone numbers (one per line)
echo "1234567890" > phones.txt

docker run -it --rm\
  -v $(pwd)/phones.txt:/app/phones.txt \
  -v $(pwd)/sessions:/app/sessions \
  ghcr.io/hnnsly/tg-sgen -app-id YOUR_APP_ID -app-hash "YOUR_APP_HASH"
```

### Local Build

```bash
go build -o sgen main.go
./sgen -app-id YOUR_APP_ID -app-hash "YOUR_APP_HASH"
```

## 📖 Usage Examples

```bash
# Basic usage
./sgen -app-id 12345678 -app-hash "your_app_hash"

# With channel joining
./sgen -app-id 12345678 -app-hash "your_hash" -channel @mychannel -verbose

# Custom paths
./sgen -app-id 12345678 -app-hash "your_hash" -phones-file phones.txt -sessions-dir sessions
```

## ⚙️ Options

| Flag            | Description                     | Default      |
| --------------- | ------------------------------- | ------------ |
| `-app-id`       | **Required.** Telegram App ID   | -            |
| `-app-hash`     | **Required.** Telegram App Hash | -            |
| `-phones-file`  | Phone numbers file              | `phones.txt` |
| `-sessions-dir` | Sessions directory              | `sessions`   |
| `-channel`      | Channel to join (repeatable)    | -            |
| `-verbose`      | Verbose logging                 | `false`      |

## 📄 phones.txt Format

```txt
1234567890
+380501234567
//1999999999, this line will be skipped
```

## 🙏 Credits

Built with [Gogram](https://github.com/amarnathcjd/gogram) by [AmarnathCJD](https://github.com/AmarnathCJD). Thanks for the amazing library!

---

💡 **Remember:** Always use your own App ID/Hash and follow Telegram ToS
