Saya akan menjelaskan lebih detail tentang perilaku konkuren untuk setiap level isolasi yang Anda tanyakan:

1. Read Committed
   Pada level ini, transaksi tidak akan menunggu transaksi lain untuk membaca data, tapi akan menunggu untuk menulis data. Mari kita lihat contohnya:

```sql
-- Transaksi 1
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
-- akan mengunci baris dengan id=1

-- Transaksi 2 (bersamaan)
BEGIN;
SELECT balance FROM accounts WHERE id = 1;
-- bisa langsung membaca data lama
-- tidak perlu menunggu Transaksi 1

UPDATE accounts SET balance = balance + 200 WHERE id = 1;
-- harus menunggu sampai Transaksi 1 selesai (COMMIT atau ROLLBACK)
-- karena baris sedang dikunci oleh Transaksi 1
```

Jadi pada Read Committed:

- Operasi SELECT tidak perlu menunggu
- Operasi UPDATE/DELETE harus menunggu jika ada transaksi lain yang sedang mengubah baris yang sama
- Setiap SELECT melihat data terbaru yang sudah di-commit

2. Repeatable Read
   Pada level ini, setiap transaksi bekerja dengan "snapshot" datanya sendiri. Mari lihat contohnya:

```sql
-- Transaksi 1
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SELECT balance FROM accounts WHERE id = 1;  -- misalnya 1000

-- Transaksi 2 (bersamaan)
BEGIN;
UPDATE accounts SET balance = 800 WHERE id = 1;
COMMIT;

-- Kembali ke Transaksi 1
SELECT balance FROM accounts WHERE id = 1;  -- tetap 1000
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
-- Akan menggunakan nilai 1000 sebagai balance, bukan 800
-- Jika perubahan dari Transaksi 2 bertentangan dengan perubahan ini
-- maka Transaksi 1 akan gagal dengan error "could not serialize access"
```

Pada Repeatable Read:

- Transaksi bisa berjalan bersamaan
- Tidak perlu menunggu untuk membaca
- Untuk UPDATE, sistem akan memeriksa apakah data yang akan diubah sudah berubah sejak transaksi dimulai
- Jika ada konflik, transaksi akan gagal dan harus diulang

3. Serializable
   Benar sekali! Pada level Serializable, PostgreSQL menggunakan teknik yang disebut "Serializable Snapshot Isolation (SSI)". Mari lihat contoh yang bisa menyebabkan deadlock:

```sql
-- Transaksi 1
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- Transaksi 2 (bersamaan)
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
UPDATE accounts SET balance = balance - 100 WHERE id = 2;
UPDATE accounts SET balance = balance + 100 WHERE id = 1;
```

Situasi di atas bisa menyebabkan deadlock karena:

- Transaksi 1 mengunci id=1 dan menunggu id=2
- Transaksi 2 mengunci id=2 dan menunggu id=1
- Keduanya saling menunggu tanpa ada yang bisa melanjutkan

Untuk menghindari deadlock pada Serializable:

1. Selalu update data dalam urutan yang sama (misal: urut berdasarkan ID)
2. Gunakan timeout pada transaksi
3. Siapkan logika untuk mencoba ulang transaksi yang gagal
4. Batasi jumlah baris yang diupdate dalam satu transaksi

Contoh penanganan yang baik:

```sql
DO $$
DECLARE
  max_retries INTEGER := 3;
  current_try INTEGER := 0;
BEGIN
  WHILE current_try < max_retries LOOP
    BEGIN
      -- Mulai transaksi
      BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;

      -- Lakukan operasi database
      UPDATE accounts SET balance = balance - 100
      WHERE id = 1;

      UPDATE accounts SET balance = balance + 100
      WHERE id = 2;

      COMMIT;
      -- Jika berhasil, keluar dari loop
      EXIT;
    EXCEPTION WHEN serialization_failure OR deadlock_detected THEN
      -- Jika gagal, tambah counter dan coba lagi
      current_try := current_try + 1;
      IF current_try = max_retries THEN
        RAISE EXCEPTION 'Gagal setelah % percobaan', max_retries;
      END IF;
    END;
  END LOOP;
END $$;
```
