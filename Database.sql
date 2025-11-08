-- =================================================================
-- 1. CONFIGURATION
-- =================================================================
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+07:00";


-- =================================================================
-- 2. DROP TABLES (In reverse order of dependency)
-- =================================================================
DROP TABLE IF EXISTS `payments`;
DROP TABLE IF EXISTS `maintenance_requests`;
DROP TABLE IF EXISTS `announcements`;
DROP TABLE IF EXISTS `leases`;
DROP TABLE IF EXISTS `monthly_billing`;
DROP TABLE IF EXISTS `rooms`;
DROP TABLE IF EXISTS `users`;


-- =================================================================
-- 3. CREATE TABLES
-- =================================================================

-- ตาราง 1: ผู้ใช้ (Users)
CREATE TABLE `users` (
  `user_id` INT(11) NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `password_salt` VARCHAR(255) NOT NULL,
  `full_name` VARCHAR(100) NOT NULL,
  `phone` VARCHAR(15) DEFAULT NULL,
  `role` ENUM('owner', 'tenant') NOT NULL DEFAULT 'tenant',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ตาราง 2: ห้องพัก (Rooms)
CREATE TABLE `rooms` (
  `room_id` INT(11) NOT NULL AUTO_INCREMENT,
  `room_number` VARCHAR(10) NOT NULL,
  `floor` INT(3) DEFAULT NULL,
  `status` ENUM('available', 'occupied', 'maintenance') NOT NULL DEFAULT 'available',
  `tenant_id` INT(11) DEFAULT NULL COMMENT 'ID ผู้เช่าปัจจุบัน (อัปเดตจาก leases)',
  `base_rent` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `furniture_fee` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  PRIMARY KEY (`room_id`),
  UNIQUE KEY `room_number` (`room_number`),
  KEY `tenant_id` (`tenant_id`),
  CONSTRAINT `rooms_ibfk_1` FOREIGN KEY (`tenant_id`) REFERENCES `users` (`user_id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ตาราง 3: บิลรายเดือน (Monthly Billing - Utilities)
CREATE TABLE `monthly_billing` (
  `billing_id` INT(11) NOT NULL AUTO_INCREMENT,
  `room_id` INT(11) NOT NULL,
  `billing_month` DATE NOT NULL COMMENT 'บิลประจำเดือน (เช่น 2025-11-01)',
  `due_date` DATE DEFAULT NULL COMMENT 'วันครบกำหนดชำระ',
  `water_units` DECIMAL(10,2) DEFAULT 0.00,
  `electricity_units` DECIMAL(10,2) DEFAULT 0.00,
  `water_bill` DECIMAL(10,2) DEFAULT 0.00,
  `electricity_bill` DECIMAL(10,2) DEFAULT 0.00,
  `total_utility_bill` DECIMAL(10,2) DEFAULT 0.00 COMMENT 'ยอดรวม (น้ำ+ไฟ)',
  `status` ENUM('unpaid', 'paid', 'overdue') NOT NULL DEFAULT 'unpaid',
  PRIMARY KEY (`billing_id`),
  KEY `room_id` (`room_id`),
  CONSTRAINT `billing_ibfk_1` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`room_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ตาราง 4: (ใหม่) สัญญาเช่า (Leases)
CREATE TABLE `leases` (
  `lease_id` INT(11) NOT NULL AUTO_INCREMENT,
  `room_id` INT(11) NOT NULL,
  `tenant_id` INT(11) NOT NULL,
  `start_date` DATE NOT NULL COMMENT 'วันเริ่มสัญญา',
  `end_date` DATE NOT NULL COMMENT 'วันสิ้นสุดสัญญา',
  `security_deposit` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT 'ค่ามัดจำ',
  `status` ENUM('active', 'expired', 'terminated') NOT NULL DEFAULT 'active',
  PRIMARY KEY (`lease_id`),
  KEY `room_id` (`room_id`),
  KEY `tenant_id` (`tenant_id`),
  CONSTRAINT `leases_ibfk_1` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`room_id`) ON UPDATE CASCADE,
  CONSTRAINT `leases_ibfk_2` FOREIGN KEY (`tenant_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ตาราง 5: (ใหม่) การชำระเงิน (Payments)
CREATE TABLE `payments` (
  `payment_id` INT(11) NOT NULL AUTO_INCREMENT,
  `billing_id` INT(11) NOT NULL COMMENT 'อ้างอิงบิลค่าน้ำค่าไฟ',
  `tenant_id` INT(11) NOT NULL,
  `amount_paid` DECIMAL(10,2) NOT NULL,
  `payment_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'วันที่และเวลาที่จ่าย',
  `payment_method` ENUM('cash', 'transfer') NOT NULL,
  `notes` VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (`payment_id`),
  KEY `billing_id` (`billing_id`),
  KEY `tenant_id` (`tenant_id`),
  CONSTRAINT `payments_ibfk_1` FOREIGN KEY (`billing_id`) REFERENCES `monthly_billing` (`billing_id`) ON UPDATE CASCADE,
  CONSTRAINT `payments_ibfk_2` FOREIGN KEY (`tenant_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ตาราง 6: (ใหม่) การแจ้งซ่อม (Maintenance Requests)
CREATE TABLE `maintenance_requests` (
  `request_id` INT(11) NOT NULL AUTO_INCREMENT,
  `room_id` INT(11) NOT NULL,
  `tenant_id` INT(11) NOT NULL,
  `issue_description` TEXT NOT NULL COMMENT 'รายละเอียดปัญหา',
  `request_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'วันที่แจ้ง',
  `status` ENUM('pending', 'in_progress', 'completed') NOT NULL DEFAULT 'pending',
  `repair_cost` DECIMAL(10,2) DEFAULT 0.00 COMMENT 'ค่าซ่อม (ถ้ามี)',
  PRIMARY KEY (`request_id`),
  KEY `room_id` (`room_id`),
  KEY `tenant_id` (`tenant_id`),
  CONSTRAINT `maint_ibfk_1` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`room_id`) ON UPDATE CASCADE,
  CONSTRAINT `maint_ibfk_2` FOREIGN KEY (`tenant_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ตาราง 7: (ใหม่) ประกาศ (Announcements)
CREATE TABLE `announcements` (
  `announcement_id` INT(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL COMMENT 'ผู้ประกาศ (เจ้าของ)',
  `title` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`announcement_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `anno_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =================================================================
-- 4. INSERT DATA
-- =================================================================

-- ข้อมูล: Users (1 เจ้าของ + 39 ผู้เช่า)
INSERT INTO `users` (`user_id`, `username`, `password_hash`, `password_salt`, `full_name`, `phone`, `role`) VALUES
(1, 'owner01', 'hash_ของเจ้าของ_xyz', 'salt_ของเจ้าของ_abc', 'สมชาย ใจดี (เจ้าของ)', '0810001111', 'owner');
INSERT INTO `users` (`username`, `password_hash`, `password_salt`, `full_name`, `phone`, `role`) VALUES
('tenant001', 'hash001', 'salt001', 'สมหญิง จริงใจ', '0810000001', 'tenant'),
('tenant002', 'hash002', 'salt002', 'สมศักดิ์ รักสงบ', '0810000002', 'tenant'),
('tenant003', 'hash003', 'salt003', 'มานี มีนา', '0810000003', 'tenant'),
('tenant004', 'hash004', 'salt004', 'ปิติ ชูใจ', '0810000004', 'tenant'),
('tenant005', 'hash005', 'salt005', 'อาทิตย์ จันทร์เพ็ญ', '0810000005', 'tenant'),
('tenant006', 'hash006', 'salt006', 'วิไล พรรณราย', '0810000006', 'tenant'),
('tenant007', 'hash007', 'salt007', 'กวินทร์ อินทร์กอง', '0810000007', 'tenant'),
('tenant008', 'hash008', 'salt008', 'นที สีดา', '0810000008', 'tenant'),
('tenant009', 'hash009', 'salt009', 'ธารา วาริน', '0810000009', 'tenant'),
('tenant010', 'hash010', 'salt010', 'เอกพล คนขยัน', '0810000010', 'tenant'),
('tenant011', 'hash011', 'salt011', 'จิราพร ออนซอน', '0810000011', 'tenant'),
('tenant012', 'hash012', 'salt012', 'ทินกร ก่อนเที่ยง', '0810000012', 'tenant'),
('tenant013', 'hash013', 'salt013', 'มรกต อัญมณี', '0810000013', 'tenant'),
('tenant014', 'hash014', 'salt014', 'เพทาย สายไหม', '0810000014', 'tenant'),
('tenant015', 'hash015', 'salt015', 'บุษราคัม คำหวาน', '0810000015', 'tenant'),
('tenant016', 'hash016', 'salt016', 'ไพลิน ถิ่นไทย', '0810000016', 'tenant'),
('tenant017', 'hash017', 'salt017', 'โกเมน เอกชัย', '0810000017', 'tenant'),
('tenant018', 'hash018', 'salt018', 'นิลกาฬ ชาตรี', '0810000018', 'tenant'),
('tenant019', 'hash019', 'salt019', 'เพชรรัตน์ มงคล', '0810000019', 'tenant'),
('tenant020', 'hash020', 'salt020', 'ทับทิม สยาม', '0810000020', 'tenant'),
('tenant021', 'hash021', 'salt021', 'วิเชียร มาศวิไล', '0810000021', 'tenant'),
('tenant022', 'hash022', 'salt022', 'สุริยา จันทรา', '0810000022', 'tenant'),
('tenant023', 'hash023', 'salt023', 'ดวงดาว พร่างพราว', '0810000023', 'tenant'),
('tenant024', 'hash024', 'salt024', 'ฟ้าใส ใจสบาย', '0810000024', 'tenant'),
('tenant025', 'hash025', 'salt025', 'อรุณ รุ่งเรือง', '0810000025', 'tenant'),
('tenant026', 'hash026', 'salt026', 'สายัณห์ ตะวันฉาย', '0810000026', 'tenant'),
('tenant027', 'hash027', 'salt027', 'ราตรี มีสุข', '0810000027', 'tenant'),
('tenant028', 'hash028', 'salt028', 'วสันต์ ฤดูดี', '0810000028', 'tenant'),
('tenant029', 'hash029', 'salt029', 'เหมันต์ วันหนาว', '0810000029', 'tenant'),
('tenant030', 'hash030', 'salt030', 'คิมหันต์ แสนร้อน', '0810000030', 'tenant'),
('tenant031', 'hash031', 'salt031', 'พิรุณ กรุณา', '0810000031', 'tenant'),
('tenant032', 'hash032', 'salt032', 'เมฆา พาฝน', '0810000032', 'tenant'),
('tenant033', 'hash033', 'salt033', 'สายลม อ่อนโยน', '0810000033', 'tenant'),
('tenant034', 'hash034', 'salt034', 'ปฐพี คงทน', '0810000034', 'tenant'),
('tenant035', 'hash035', 'salt035', 'อัคคี ร้อนแรง', '0810000035', 'tenant'),
('tenant036', 'hash036', 'salt036', 'วายุ พัดพา', '0810000036', 'tenant'),
('tenant037', 'hash037', 'salt037', 'ชลธี ไหลริน', '0810000037', 'tenant'),
('tenant038', 'hash038', 'salt038', 'พนาไพร เขียวขจี', '0810000038', 'tenant'),
('tenant039', 'hash039', 'salt039', 'ภูผา สูงใหญ่', '0810000039', 'tenant');

-- ข้อมูล: Rooms (40 ห้อง)
-- (ผูกผู้เช่า user_id 2 ถึง 40 เข้ากับห้อง room_id 1 ถึง 39)
INSERT INTO `rooms` (`room_id`, `room_number`, `floor`, `status`, `tenant_id`, `base_rent`, `furniture_fee`) VALUES
(1, '101', 1, 'occupied', 2, 4500.00, 500.00),
(2, '102', 1, 'occupied', 3, 4500.00, 500.00),
(3, '103', 1, 'occupied', 4, 4500.00, 500.00),
(4, '104', 1, 'occupied', 5, 4500.00, 500.00),
(5, '105', 1, 'occupied', 6, 4500.00, 0.00),
(6, '201', 2, 'occupied', 7, 5000.00, 500.00),
(7, '202', 2, 'occupied', 8, 5000.00, 500.00),
(8, '203', 2, 'occupied', 9, 5000.00, 500.00),
(9, '204', 2, 'occupied', 10, 5000.00, 0.00),
(10, '205', 2, 'occupied', 11, 5000.00, 0.00),
(11, '301', 3, 'occupied', 12, 5500.00, 500.00),
(12, '302', 3, 'occupied', 13, 5500.00, 500.00),
(13, '303', 3, 'occupied', 14, 5500.00, 500.00),
(14, '304', 3, 'occupied', 15, 5500.00, 0.00),
(15, '305', 3, 'occupied', 16, 5500.00, 0.00),
(16, '401', 4, 'occupied', 17, 6000.00, 500.00),
(17, '402', 4, 'occupied', 18, 6000.00, 500.00),
(18, '403', 4, 'occupied', 19, 6000.00, 500.00),
(19, '404', 4, 'occupied', 20, 6000.00, 0.00),
(20, '405', 4, 'occupied', 21, 6000.00, 0.00),
(21, '501', 5, 'occupied', 22, 6500.00, 500.00),
(22, '502', 5, 'occupied', 23, 6500.00, 500.00),
(23, '503', 5, 'occupied', 24, 6500.00, 500.00),
(24, '504', 5, 'occupied', 25, 6500.00, 0.00),
(25, '505', 5, 'occupied', 26, 6500.00, 0.00),
(26, '601', 6, 'occupied', 27, 7000.00, 500.00),
(27, '602', 6, 'occupied', 28, 7000.00, 500.00),
(28, '603', 6, 'occupied', 29, 7000.00, 500.00),
(29, '604', 6, 'occupied', 30, 7000.00, 0.00),
(30, '605', 6, 'occupied', 31, 7000.00, 0.00),
(31, '701', 7, 'occupied', 32, 7500.00, 500.00),
(32, '702', 7, 'occupied', 33, 7500.00, 500.00),
(33, '703', 7, 'occupied', 34, 7500.00, 500.00),
(34, '704', 7, 'occupied', 35, 7500.00, 0.00),
(35, '705', 7, 'occupied', 36, 7500.00, 0.00),
(36, '801', 8, 'occupied', 37, 8000.00, 500.00),
(37, '802', 8, 'occupied', 38, 8000.00, 500.00),
(38, '803', 8, 'occupied', 39, 8000.00, 500.00),
(39, '804', 8, 'occupied', 40, 8000.00, 0.00),
(40, '805', 8, 'available', NULL, 8000.00, 0.00);

-- ข้อมูล: Monthly Billing (ครบ 39 บิล)
INSERT INTO `monthly_billing` (`billing_id`, `room_id`, `billing_month`, `due_date`, `water_units`, `electricity_units`, `water_bill`, `electricity_bill`, `total_utility_bill`, `status`) VALUES
(1, 1, '2025-11-01', '2025-11-05', 10.0, 120.0, 180.00, 840.00, 1020.00, 'overdue'),
(2, 2, '2025-11-01', '2025-11-30', 8.0, 90.0, 144.00, 630.00, 774.00, 'paid'),
(3, 3, '2025-11-01', '2025-11-30', 12.0, 150.0, 216.00, 1050.00, 1266.00, 'unpaid'),
(4, 4, '2025-11-01', '2025-11-30', 9.0, 110.0, 162.00, 770.00, 932.00, 'unpaid'),
(5, 5, '2025-11-01', '2025-11-05', 15.0, 200.0, 270.00, 1400.00, 1670.00, 'overdue'),
(6, 6, '2025-11-01', '2025-11-30', 11.0, 130.0, 198.00, 910.00, 1108.00, 'paid'),
(7, 7, '2025-11-01', '2025-11-30', 7.0, 80.0, 126.00, 560.00, 686.00, 'paid'),
(8, 8, '2025-11-01', '2025-11-30', 10.0, 100.0, 180.00, 700.00, 880.00, 'unpaid'),
(9, 9, '2025-11-01', '2025-11-30', 13.0, 160.0, 234.00, 1120.00, 1354.00, 'unpaid'),
(10, 10, '2025-11-01', '2025-11-05', 14.0, 180.0, 252.00, 1260.00, 1512.00, 'overdue'),
(11, 11, '2025-11-01', '2025-11-30', 10.0, 120.0, 180.00, 840.00, 1020.00, 'paid'),
(12, 12, '2025-11-01', '2025-11-30', 8.0, 90.0, 144.00, 630.00, 774.00, 'paid'),
(13, 13, '2025-11-01', '2025-11-30', 12.0, 150.0, 216.00, 1050.00, 1266.00, 'unpaid'),
(14, 14, '2025-11-01', '2025-11-30', 9.0, 110.0, 162.00, 770.00, 932.00, 'unpaid'),
(15, 15, '2025-11-01', '2025-11-30', 15.0, 200.0, 270.00, 1400.00, 1670.00, 'paid'),
(16, 16, '2025-11-01', '2025-11-30', 11.0, 130.0, 198.00, 910.00, 1108.00, 'paid'),
(17, 17, '2025-11-01', '2025-11-30', 7.0, 80.0, 126.00, 560.00, 686.00, 'unpaid'),
(18, 18, '2025-11-01', '2025-11-05', 10.0, 100.0, 180.00, 700.00, 880.00, 'overdue'),
(19, 19, '2025-11-01', '2025-11-30', 13.0, 160.0, 234.00, 1120.00, 1354.00, 'paid'),
(20, 20, '2025-11-01', '2025-11-30', 14.0, 180.0, 252.00, 1260.00, 1512.00, 'unpaid'),
(21, 21, '2025-11-01', '2025-11-30', 10.0, 120.0, 180.00, 840.00, 1020.00, 'unpaid'),
(22, 22, '2025-11-01', '2025-11-05', 8.0, 90.0, 144.00, 630.00, 774.00, 'overdue'),
(23, 23, '2025-11-01', '2025-11-30', 12.0, 150.0, 216.00, 1050.00, 1266.00, 'paid'),
(24, 24, '2025-11-01', '2025-11-30', 9.0, 110.0, 162.00, 770.00, 932.00, 'unpaid'),
(25, 25, '2025-11-01', '2025-11-30', 15.0, 200.0, 270.00, 1400.00, 1670.00, 'paid'),
(26, 26, '2025-11-01', '2025-11-30', 11.0, 130.0, 198.00, 910.00, 1108.00, 'unpaid'),
(27, 27, '2025-11-01', '2025-11-30', 7.0, 80.0, 126.00, 560.00, 686.00, 'paid'),
(28, 28, '2025-11-01', '2025-11-30', 10.0, 100.0, 180.00, 700.00, 880.00, 'unpaid'),
(29, 29, '2025-11-01', '2025-11-05', 13.0, 160.0, 234.00, 1120.00, 1354.00, 'overdue'),
(30, 30, '2025-11-01', '2025-11-30', 14.0, 180.0, 252.00, 1260.00, 1512.00, 'paid'),
(31, 31, '2025-11-01', '2025-11-30', 10.0, 120.0, 180.00, 840.00, 1020.00, 'unpaid'),
(32, 32, '2025-11-01', '2025-11-30', 8.0, 90.0, 144.00, 630.00, 774.00, 'paid'),
(33, 33, '2025-11-01', '2025-11-30', 12.0, 150.0, 216.00, 1050.00, 1266.00, 'unpaid'),
(34, 34, '2025-11-01', '2025-11-05', 9.0, 110.0, 162.00, 770.00, 932.00, 'overdue'),
(35, 35, '2025-11-01', '2025-11-30', 15.0, 200.0, 270.00, 1400.00, 1670.00, 'paid'),
(36, 36, '2025-11-01', '2025-11-30', 11.0, 130.0, 198.00, 910.00, 1108.00, 'unpaid'),
(37, 37, '2025-11-01', '2025-11-30', 7.0, 80.0, 126.00, 560.00, 686.00, 'paid'),
(38, 38, '2025-11-01', '2025-11-30', 10.0, 100.0, 180.00, 700.00, 880.00, 'unpaid'),
(39, 39, '2025-11-01', '2025-11-05', 13.0, 160.0, 234.00, 1120.00, 1354.00, 'overdue');

-- ข้อมูล: Leases (39 สัญญา)
-- (สมมติว่าทุกคนมีสัญญา 1 ปี และจ่ายค่ามัดจำ = 2 * ค่าเช่าห้อง)
INSERT INTO `leases` (`lease_id`, `room_id`, `tenant_id`, `start_date`, `end_date`, `security_deposit`, `status`) VALUES
(1, 1, 2, '2025-01-01', '2025-12-31', 9000.00, 'active'),
(2, 2, 3, '2025-01-01', '2025-12-31', 9000.00, 'active'),
(3, 3, 4, '2025-01-01', '2025-12-31', 9000.00, 'active'),
(4, 4, 5, '2025-01-01', '2025-12-31', 9000.00, 'active'),
(5, 5, 6, '2025-01-01', '2025-12-31', 9000.00, 'active'),
(6, 6, 7, '2025-02-01', '2026-01-31', 10000.00, 'active'),
(7, 7, 8, '2025-02-01', '2026-01-31', 10000.00, 'active'),
(8, 8, 9, '2025-02-01', '2026-01-31', 10000.00, 'active'),
(9, 9, 10, '2025-02-01', '2026-01-31', 10000.00, 'active'),
(10, 10, 11, '2025-02-01', '2026-01-31', 10000.00, 'active'),
(11, 11, 12, '2025-03-01', '2026-02-28', 11000.00, 'active'),
(12, 12, 13, '2025-03-01', '2026-02-28', 11000.00, 'active'),
(13, 13, 14, '2025-03-01', '2026-02-28', 11000.00, 'active'),
(14, 14, 15, '2025-03-01', '2026-02-28', 11000.00, 'active'),
(15, 15, 16, '2025-03-01', '2026-02-28', 11000.00, 'active'),
(16, 16, 17, '2025-04-01', '2026-03-31', 12000.00, 'active'),
(17, 17, 18, '2025-04-01', '2026-03-31', 12000.00, 'active'),
(18, 18, 19, '2025-04-01', '2026-03-31', 12000.00, 'active'),
(19, 19, 20, '2025-04-01', '2026-03-31', 12000.00, 'active'),
(20, 20, 21, '2025-04-01', '2026-03-31', 12000.00, 'active'),
(21, 21, 22, '2025-05-01', '2026-04-30', 13000.00, 'active'),
(22, 22, 23, '2025-05-01', '2026-04-30', 13000.00, 'active'),
(23, 23, 24, '2025-05-01', '2026-04-30', 13000.00, 'active'),
(24, 24, 25, '2025-05-01', '2026-04-30', 13000.00, 'active'),
(25, 25, 26, '2025-05-01', '2026-04-30', 13000.00, 'active'),
(26, 26, 27, '2025-06-01', '2026-05-31', 14000.00, 'active'),
(27, 27, 28, '2025-06-01', '2026-05-31', 14000.00, 'active'),
(28, 28, 29, '2025-06-01', '2026-05-31', 14000.00, 'active'),
(29, 29, 30, '2025-06-01', '2026-05-31', 14000.00, 'active'),
(30, 30, 31, '2025-06-01', '2026-05-31', 14000.00, 'active'),
(31, 31, 32, '2025-07-01', '2026-06-30', 15000.00, 'active'),
(32, 32, 33, '2025-07-01', '2026-06-30', 15000.00, 'active'),
(33, 33, 34, '2025-07-01', '2026-06-30', 15000.00, 'active'),
(34, 34, 35, '2025-07-01', '2026-06-30', 15000.00, 'active'),
(35, 35, 36, '2025-07-01', '2026-06-30', 15000.00, 'active'),
(36, 36, 37, '2025-08-01', '2026-07-31', 16000.00, 'active'),
(37, 37, 38, '2025-08-01', '2026-07-31', 16000.00, 'active'),
(38, 38, 39, '2025-08-01', '2026-07-31', 16000.00, 'active'),
(39, 39, 40, '2025-08-01', '2026-07-31', 16000.00, 'active');

-- ข้อมูล: Payments (15 รายการ, สำหรับบิลที่ 'paid')
-- (สมมติว่าจ่ายค่าน้ำค่าไฟ ตามยอด total_utility_bill)
INSERT INTO `payments` (`payment_id`, `billing_id`, `tenant_id`, `amount_paid`, `payment_date`, `payment_method`, `notes`) VALUES
(1, 2, 3, 774.00, '2025-11-02 10:30:00', 'transfer', 'โอนจาก KBank'),
(2, 6, 7, 1108.00, '2025-11-02 11:15:00', 'cash', 'จ่ายเงินสดที่ออฟฟิศ'),
(3, 7, 8, 686.00, '2025-11-03 09:00:00', 'transfer', 'โอนจาก SCB'),
(4, 11, 12, 1020.00, '2025-11-03 14:22:00', 'transfer', 'โอน'),
(5, 12, 13, 774.00, '2025-11-03 15:00:00', 'cash', 'เงินสด'),
(6, 15, 16, 1670.00, '2025-11-04 08:10:00', 'transfer', 'โอนแล้ว'),
(7, 16, 17, 1108.00, '2025-11-04 12:00:00', 'cash', 'จ่ายสด'),
(8, 19, 20, 1354.00, '2025-11-04 13:00:00', 'transfer', NULL),
(9, 23, 24, 1266.00, '2025-11-05 10:00:00', 'transfer', 'KBank'),
(10, 25, 26, 1670.00, '2025-11-05 11:00:00', 'cash', NULL),
(11, 27, 28, 686.00, '2025-11-06 16:00:00', 'transfer', 'โอน'),
(12, 30, 31, 1512.00, '2025-11-06 17:00:00', 'transfer', NULL),
(13, 32, 33, 774.00, '2025-11-07 10:00:00', 'cash', 'ฝากไว้ที่เคาน์เตอร์'),
(14, 35, 36, 1670.00, '2025-11-07 11:00:00', 'transfer', 'SCB'),
(15, 37, 38, 686.00, '2025-11-08 13:00:00', 'transfer', 'เรียบร้อย');

-- ข้อมูล: Maintenance Requests (4 รายการ)
INSERT INTO `maintenance_requests` (`request_id`, `room_id`, `tenant_id`, `issue_description`, `request_date`, `status`, `repair_cost`) VALUES
(1, 1, 2, 'ก๊อกน้ำรั่วซึมที่อ่างล้างหน้า','2025-11-01 09:00:00', 'pending', 0.00),
(2, 6, 7, 'แอร์ไม่เย็นเลยครับ', '2025-11-03 14:00:00', 'in_progress', 0.00),
(3, 11, 12, 'หลอดไฟในห้องน้ำขาด', '2025-11-04 10:00:00', 'completed', 150.00),
(4, 22, 23, 'ชักโครกกดไม่ลง', '2025-11-05 18:00:00', 'pending', 0.00);

-- ข้อมูล: Announcements (2 รายการ)
INSERT INTO `announcements` (`announcement_id`, `user_id`, `title`, `content`, `created_at`) VALUES
(1, 1, 'แจ้งกำหนดการตัดไฟ 10 พ.ย. 2568', 'เนื่องด้วยการไฟฟ้าจะทำการปรับปรุงหม้อแปลง หอพักจะไฟดับในวันที่ 10 พฤศจิกายน 2568 เวลา 09:00 - 15:00 น. ขออภัยในความไม่สะดวก', '2025-11-08 10:00:00'),
(2, 1, 'ประกาศ: ห้ามเลี้ยงสัตว์ทุกชนิด', 'เน้นย้ำเรื่องกฎระเบียบ ห้ามมิให้ผู้เช่านำสัตว์เลี้ยงทุกชนิดเข้ามาเลี้ยงในบริเวณหอพักโดยเด็ดขาด หากตรวจพบจะดำเนินการตามสัญญาเช่า', '2025-11-09 08:00:00');


-- =================================================================
-- 5. COMMIT
-- =================================================================
COMMIT;