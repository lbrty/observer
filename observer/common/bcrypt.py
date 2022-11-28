import bcrypt
from password_strength import PasswordPolicy

"""
+------+-------------------+
| Cost | Iterations        |
|------|-------------------|
|  8   |    256 iterations |
|  9   |    512 iterations |
| 10   |  1,024 iterations |
| 11   |  2,048 iterations |
| 12   |  4,096 iterations |
| 13   |  8,192 iterations |
| 14   | 16,384 iterations |
| 15   | 32,768 iterations |
| 16   | 65,536 iterations |
+------+-------------------+
"""
HashingRounds = 12


def check_password(password: str, password_hash: str) -> bool:
    return bcrypt.checkpw(password.encode(), password_hash.encode())


def hash_password(password: str) -> str:
    return bcrypt.hashpw(password.encode(), bcrypt.gensalt(rounds=HashingRounds)).decode()


def is_strong_password(password: str, policy: PasswordPolicy) -> bool:
    return len(policy.test(password)) == 0
