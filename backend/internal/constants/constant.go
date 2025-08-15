package constants

const (
	// ContentType is key for get reponse content type in context
	ContentType = "Content-Type"

	Success                 = "00"
	InternalErr             = "99"
	SystemError             = "01" //
	TransactionNotProcess   = "02" //02 Transaksi Tidak Dapat Diproses
	CommunicationError      = "03" //03 Communication Error
	TransactionNotValid     = "04" //04 Transaksi Tidak Valid
	TransactionSuspect      = "05" //05 Transaksi Suspect
	Timeout                 = "06" //06 Timeout
	TerminalNotValid        = "07" // 07 Terminal ID Tidak Valid
	MerchantNotValid        = "08" //08 Merchant ID Tidak Valid
	AmountNotValid          = "09" //09 Amount Transaksi Tidak Valid Atau Duplikasi
	CustomerNotRegisered    = "10" //10 Nasabah Belum Terdaftar Di eBanking
	CustomerIsRegistered    = "11" //11 Nasabah Sudah Terdaftar Di eBanking
	RegistrationFailed      = "12" //12 Registrasi Gagal
	InvalidMessage          = "13" //13 Format Message Tidak Sesuai
	UndefinedError          = "14" //14 Kode Kesalahan Tidak Terdefinisi
	DatabaseConnError       = "15" //15 Koneksi Ke Database Gagal
	AmountLessThanLimit     = "16" // Nilai Transaksi Kurang Dari Batas Minimum
	InvalidInstitution      = "17" // Intitusi Tidak Valid
	TransactionNotFound     = "18" //Nomor Referensi / Data Transaksi Tidak Ditemukan
	SystemOffline           = "19" //Sistem Sedang Offline
	DifferentSettlementDate = "20" //Tanggal Settlement Berbeda
	PrevTrxBeingProcc       = "21" //Transaksi Sebelumnya Sedang Diproses
	TransactionNotAvailable = "22" // Transaksi tidak tersedia|Transaction is not available
	MaxReprintExeed         = "23" // Cetak ulang mencapai limit|Maxium reprint exceed
	AmountMoreThanLimit     = "24" // Nilai Transaksi Melebihi Batas Maximum

	//Customer
	CustomerNotFound = "30" //Nasabah Tidak Ditemukan|Customer Not Found

	//ACCOUNT, CARD & LIMIT
	InvalidAccNumber   = "41"
	AccNumberNotActive = "42"
	AccNumberBlocked   = "43"
	InvalidBalance     = "44"

	//Biller

	InvalidQR = "6G" // 6G	Kode QR Tidak Valid|Invalid QR Code

	//Authentication
	InvalidAuthentication      = "70" //Otentikasi gagal|Invalid authentication
	UsernameNotFound           = "71" //Username tidak terdaftar|Username is not found
	UsernameNotActive          = "72" //Username tidak aktif|Username is not active
	CustomerNotActive          = "73" //Nasabah tidak aktif|Customer is not active
	DeviceNotActive            = "74" //Gadget tidak aktif|Device is not active
	InvalidOldPassword         = "75" //Password lama tidak benar|Invalid old password
	OldAndNewPassMustDifferent = "76" //Password lama dan baru tidak boleh sama|old and new password should be different
	InvalidOldPin              = "77" // PIN lama tidak benar|Invalid old PIN
	OldAndNewPINMustDifferent  = "78" // PIN lama dan baru tidak boleh sama|Old and new PIN should be different
	InvalidCustomer            = "79" // Nasabah tidak valid|Invalid customer
	InvalidOTP                 = "7A" //OTP tidak benar|Invalid OTP
	InvalidToken               = "7B" //Token tidak benar|Invalid Token
	InvalidVoucherCode         = "7C" //Kode voucher tidak benar|Invalid voucher code
	ExpireVoucherCode          = "7D" // Kode voucher kadaluarsa|Voucher Code Expired
	DeviceNotRegistered        = "7E" // Device ID tidak terdaftar|Device ID is not registered
	UserBlocked                = "7F" // User di blok|User is blocked
	DuplicateSession           = "7G" //Session Ganda|Duplicate Session
	MaxAttmeptPassword         = "7H" //Kesalahan input password lama melebihi batas, akun anda akan di blok|Maximum old password attempt limit reached,, your account is blocked
	MaxAttemptPin              = "7I" // Kesalahan input mpin lama melebihi batas akun anda akan di blok|Maximum old mpin attempt limit reached, your account is blocked

	//General
	DataNotFound                   = "80"
	InvalidRequest                 = "81" //Request tidak valid|Invalid request data
	InvalidAmount                  = "82" //Nominal salah|Invalid amount
	InvalidFee                     = "83" //Biaya admin salah|Invalid admin fee
	InvalidDate                    = "84" //Kesalahan waktu transaksi|Invalid date time
	DatabaseError                  = "85" //Kesalahan database|Database Error
	TrfListRegistered              = "86" // Transfer list sudah terdaftar
	PrefixNotFound                 = "87" // 87 Prefix Operator ditemukan|Cellular prefix number not found
	AccStatementNotFound           = "8A" //Mutasi Rekening Tidak Ditemukan|Account statement not found
	TransferListRegisterNotAllowed = "8C" //Tidak boleh menambahkan tujuan rekening yang digunakan pada saat registrasi|Transfer list registration is not allowed for this account

)

const (

	//General
	ResSuccess          = "Sukses|Success"
	ResInternalError    = "Kesalahan Internal Server|Internal Server Error"
	ResSystemError      = "Sistem Mengalami Gangguan"
	ResTrxCannotProcess = "Transaksi Tidak Dapat Diproses"

	//Customer
	ResCustomerNotFound = "Nasabah Tidak Ditemukan|Customer Not Found"

	//Authentication
	ResInvalidAuthentication      = "Otentikasi gagal|Invalid authentication"
	ResUsernameNotFound           = "Username tidak terdaftar|Username is not found"
	ResUsernameNotActive          = "Username tidak aktif|Username is not active"
	ResCustomerNotActive          = "Nasabah tidak aktif|Customer is not active"
	ResDeviceNotActive            = "Gadget tidak aktif|Device is not active"
	ResInvalidOldPassword         = "Password lama tidak benar|Invalid old password"
	ResOldAndNewPassMustDifferent = "Password lama dan baru tidak boleh sama|old and new password should be different"
	ResInvalidOldPin              = "PIN lama tidak benar|Invalid old PIN"
	ResOldAndNewPINMustDifferent  = "PIN lama dan baru tidak boleh sama|Old and new PIN should be different"
	ResInvalidCustomer            = "Nasabah tidak valid|Invalid customer"
	ResInvalidOTP                 = "OTP tidak benar|Invalid OTP"
	ResInvalidToken               = "Token tidak benar|Invalid Token"
	ResInvalidVoucherCode         = "Kode voucher tidak benar|Invalid voucher code"
	ResExpireVoucherCode          = "Kode voucher kadaluarsa|Voucher Code Expired"
	ResDeviceNotRegistered        = "Device ID tidak terdaftar|Device ID is not registered"
	ResUserBlocked                = "User di blok|User is blocked"
	ResDuplicateSession           = "Session Ganda|Duplicate Session"
	ResMaxAttmeptPassword         = "Kesalahan input password lama melebihi batas, akun anda akan di blok|Maximum old password attempt limit reached,, your account is blocked"
	ResMaxAttemptPin              = "Kesalahan input mpin lama melebihi batas akun anda akan di blok|Maximum old mpin attempt limit reached, your account is blocked"

	ResDataNotFound                   = "Data tidak ditemukan|Data not found"
	ResInvalidRequest                 = "Request tidak valid|Invalid request data"
	ResInvalidAmount                  = "Nominal salah|Invalid amount"
	ResInvalidFee                     = "Biaya admin salah|Invalid admin fee"
	ResInvalidDate                    = "Kesalahan waktu transaksi|Invalid date time"
	ResDatabaseError                  = "Kesalahan database|Database Error"
	ResTrfListRegistered              = "Transfer list sudah terdaftar"
	ResPrefixNotFound                 = "Prefix Operator tidak ditemukan|Cellular prefix number not found"
	ResAccStatementNotFound           = "Mutasi Rekening Tidak Ditemukan|Account statement not found"
	ResTransferListRegisterNotAllowed = "Tidak boleh menambahkan tujuan rekening yang digunakan pada saat registrasi|Transfer list registration is not allowed for this account"
	ResponseTimeout                   = "Waktu habis|Timeout"
)
