# Writing Good Tests

Gunakan referensi ini ketika menulis atau mengubah test, menambahkan mock, atau membuat helper test.

## Prinsip

1. **Setiap test menyebutkan break yang ditangkapnya.** Sebelum menulis body, jelaskan perubahan production apa yang seharusnya membuat test gagal.
2. **Setiap test menjalankan behavior nyata.** Uji hasil, side effect, atau contract di boundary; jangan hanya membuktikan mock dipanggil.
3. **Expected value diturunkan secara independen.** Gunakan literal atau fixture yang diperiksa manual, bukan helper atau builder yang sama dengan production code.
4. **Test behavior, bukan source text.** Untuk script, skill, atau konfigurasi, jalankan consumer atau boundary-nya dan periksa output, side effect, atau exit code.

## Mock dan Helper

- Jangan menambahkan mock sebelum memahami side effect dependency nyata.
- Mock hanya level operasi yang lambat atau eksternal; biarkan bagian yang menjadi contract tetap nyata.
- Mock harus meniru struktur data nyata secara lengkap.
- Jangan assert keberadaan mock atau test ID mock.
- Helper yang hanya diperlukan test harus tinggal di test utility, bukan production class.
- Jika setup mock lebih rumit daripada behavior yang diuji, pertimbangkan integration test dengan komponen nyata.

## Bentuk Test

- Satu test menguji satu behavior.
- Nama test menjelaskan hasil atau kontrak yang diharapkan.
- Gunakan fixture literal untuk success, error, malformed, empty, default, unauthorized, dan boundary input bila relevan.
- Test harus dapat gagal karena branch salah, side effect hilang, argumen salah, return kosong, atau validasi hilang.

## Mutation Check

Sebelum target ditandai selesai, pikirkan mutasi realistis berikut dan pastikan setidaknya satu test menangkapnya:

- constant atau argument salah;
- handler atau branch salah;
- state change atau side effect dihilangkan;
- return dibuat kosong;
- validasi input kosong, nol, nil, unauthorized, atau malformed dihilangkan.

## Red Flags

- Setup dan assertion menggunakan object/builder yang sama sehingga selalu sama.
- Expected value dihitung oleh kode yang sedang diuji.
- Test hanya gagal karena crash atau selector hilang.
- Test meng-grep source text alih-alih menguji efeknya.
- Assertion hanya memeriksa mock.
- Mock dibuat “untuk berjaga-jaga” tanpa alasan side effect.
- Setup mock memenuhi lebih dari separuh isi test.
