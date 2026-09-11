#!/usr/bin/env python3
# fake_telnetd.py — dispositivo IoT falso para o lab
# aceita root:xc3511 e outras credenciais Mirai, responde ECCHI token
import socket, threading, sys, time

BIND = "0.0.0.0"
PORT = int(sys.argv[1]) if len(sys.argv) > 1 else 23

AUTH_TABLE = {
    "root": ["xc3511", "vizxv", "12345", "123456", "icatch99", "Zte521"],
    "admin": ["admin", "password", "1234", "123456"],
    "ffadmin": ["ffadminff"],
    "unipi": ["unipi.technology"],
    "supervisor": ["zyad1234"],
}

TOKEN_QUERY = "/bin/busybox ECCHI"
TOKEN_RESPONSE = "ECCHI: applet not found"

def handle(conn, addr):
    try:
        conn.sendall(b"\xff\xfd\x1f\xff\xfd\x20\xff\xfd\x18\xff\xfd\x23\xff\xfd\x27")
        time.sleep(0.1)
        conn.sendall(b"login: ")
        user = conn.recv(256).decode(errors="ignore").strip()
        conn.sendall(b"password: ")
        password = conn.recv(256).decode(errors="ignore").strip()

        if user in AUTH_TABLE and password in AUTH_TABLE[user]:
            conn.sendall(b"\r\n# ")
        else:
            conn.sendall(b"Login incorrect\r\n")
            conn.close()
            return

        while True:
            data = conn.recv(4096)
            if not data:
                break
            cmd = data.decode(errors="ignore").strip()
            if TOKEN_QUERY in cmd:
                conn.sendall((TOKEN_RESPONSE + "\r\n# ").encode())
            elif "wget" in cmd or "curl" in cmd:
                conn.sendall(b"Connecting...\r\nSaving to: /tmp/b\r\n# ")
            elif "chmod" in cmd:
                conn.sendall(b"# ")
            elif "exit" in cmd:
                break
            else:
                conn.sendall(b"sh: applet not found\r\n# ")
    except Exception:
        pass
    finally:
        conn.close()

def main():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind((BIND, PORT))
    s.listen(64)
    print(f"[fake_telnetd] listening on {BIND}:{PORT}")
    while True:
        conn, addr = s.accept()
        threading.Thread(target=handle, args=(conn, addr), daemon=True).start()

if __name__ == "__main__":
    main()
