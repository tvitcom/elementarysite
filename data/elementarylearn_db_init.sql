-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Aug 19, 2021 at 04:37 PM
-- Server version: 8.0.26-0ubuntu0.20.04.2
-- PHP Version: 7.4.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `elementarylearn`
--

-- --------------------------------------------------------

--
-- Table structure for table `device`
--

CREATE TABLE `device` (
  `device_id` bigint UNSIGNED NOT NULL,
  `session_id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `user_id` bigint UNSIGNED NOT NULL DEFAULT '0',
  `ua` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `screen` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `langs` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `platform` enum('pc','notebook','pad','phone','watch','tv','auto','gamepad','etc') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'etc' COMMENT 'pc, notebook, pad, phone, watch, tv, auto, gamepad, etc',
  `created` timestamp NOT NULL,
  `updated` timestamp NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='users devices here';

-- --------------------------------------------------------

--
-- Table structure for table `invoice`
--

CREATE TABLE `invoice` (
  `invoice_id` bigint UNSIGNED NOT NULL,
  `creator_id` bigint UNSIGNED NOT NULL,
  `acceptor_id` bigint UNSIGNED NOT NULL COMMENT 'user_id who accept money from user_id',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `appointed` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `mediacontent`
--

CREATE TABLE `mediacontent` (
  `mediacontent_id` bigint UNSIGNED NOT NULL,
  `event_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `filename` varchar(127) NOT NULL,
  `description` varchar(127) NOT NULL,
  `mediatype` tinyint(1) NOT NULL,
  `storage_host` varchar(45) DEFAULT '',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `moderator_id` bigint UNSIGNED DEFAULT '0' COMMENT 'if modareted then wirte user_id',
  `uploader_useragent` char(255) NOT NULL,
  `uploader_ip` char(128) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='describe type and path file of exists content';

-- --------------------------------------------------------

--
-- Table structure for table `order`
--

CREATE TABLE `order` (
  `order_id` bigint UNSIGNED NOT NULL,
  `invoice_id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED NOT NULL,
  `qty` int NOT NULL,
  `options` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `payment`
--

CREATE TABLE `payment` (
  `payment_id` bigint UNSIGNED NOT NULL,
  `payer_id` bigint UNSIGNED NOT NULL,
  `acceptor_id` bigint UNSIGNED NOT NULL COMMENT 'user_id - money for',
  `tm_when` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `payment_notes` varchar(128) NOT NULL,
  `money` int NOT NULL,
  `currency_iso` char(3) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `product`
--

CREATE TABLE `product` (
  `product_id` bigint UNSIGNED NOT NULL,
  `product_category_id` bigint UNSIGNED NOT NULL,
  `title` varchar(45) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `price` int NOT NULL DEFAULT '0',
  `currency_iso` char(3) NOT NULL DEFAULT 'UAH',
  `peritem` varchar(45) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `product_category`
--

CREATE TABLE `product_category` (
  `product_category_id` bigint UNSIGNED NOT NULL,
  `alias` varchar(12) DEFAULT NULL,
  `name` varchar(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `product_detail`
--

CREATE TABLE `product_detail` (
  `product_detail_id` bigint UNSIGNED NOT NULL,
  `product_detail_group_id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `keywords` varchar(128) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `product_detail_group`
--

CREATE TABLE `product_detail_group` (
  `product_detail_group_id` bigint UNSIGNED NOT NULL,
  `name` varchar(45) NOT NULL,
  `keywords` varchar(128) NOT NULL DEFAULT '',
  `description` varchar(128) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `rating`
--

CREATE TABLE `rating` (
  `rating_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `item_kind` enum('users','quotes','reviews','mediacontent') NOT NULL,
  `item_id` bigint UNSIGNED NOT NULL,
  `rate` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `review`
--

CREATE TABLE `review` (
  `review_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `msg` varchar(150) NOT NULL,
  `photo_filename` varchar(128) NOT NULL DEFAULT '',
  `moderator_id` bigint UNSIGNED DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='reply from registered user';

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `user_id` bigint UNSIGNED NOT NULL,
  `name` varchar(64) NOT NULL,
  `email` varchar(64) NOT NULL,
  `tel` varchar(15) NOT NULL,
  `impp` varchar(45) NOT NULL DEFAULT '',
  `auth_key` varchar(128) NOT NULL DEFAULT '',
  `password_hash` varchar(128) NOT NULL,
  `approve_token` varchar(128) NOT NULL DEFAULT '',
  `created` timestamp NOT NULL,
  `lastlogin` timestamp NOT NULL,
  `roles` varchar(32) NOT NULL DEFAULT 'guest',
  `notes` varchar(128) NOT NULL DEFAULT '',
  `picture` varchar(128) NOT NULL DEFAULT '',
  `address` varchar(95) NOT NULL DEFAULT '',
  `tos_agree` varchar(64) NOT NULL DEFAULT '' COMMENT 'Terms of service agreement string Format: ver:dt:platform:ip',
  `privacy_agree` varchar(64) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `device`
--
ALTER TABLE `device`
  ADD PRIMARY KEY (`device_id`);

--
-- Indexes for table `invoice`
--
ALTER TABLE `invoice`
  ADD PRIMARY KEY (`invoice_id`);

--
-- Indexes for table `mediacontent`
--
ALTER TABLE `mediacontent`
  ADD PRIMARY KEY (`mediacontent_id`);

--
-- Indexes for table `order`
--
ALTER TABLE `order`
  ADD PRIMARY KEY (`order_id`),
  ADD KEY `fk_order_invoice1_idx` (`invoice_id`);

--
-- Indexes for table `payment`
--
ALTER TABLE `payment`
  ADD PRIMARY KEY (`payment_id`);

--
-- Indexes for table `product`
--
ALTER TABLE `product`
  ADD PRIMARY KEY (`product_id`);

--
-- Indexes for table `product_category`
--
ALTER TABLE `product_category`
  ADD PRIMARY KEY (`product_category_id`);

--
-- Indexes for table `product_detail`
--
ALTER TABLE `product_detail`
  ADD PRIMARY KEY (`product_detail_id`);

--
-- Indexes for table `product_detail_group`
--
ALTER TABLE `product_detail_group`
  ADD PRIMARY KEY (`product_detail_group_id`);

--
-- Indexes for table `rating`
--
ALTER TABLE `rating`
  ADD PRIMARY KEY (`rating_id`);

--
-- Indexes for table `review`
--
ALTER TABLE `review`
  ADD PRIMARY KEY (`review_id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `email_UNIQUE` (`email`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `device`
--
ALTER TABLE `device`
  MODIFY `device_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `invoice`
--
ALTER TABLE `invoice`
  MODIFY `invoice_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `mediacontent`
--
ALTER TABLE `mediacontent`
  MODIFY `mediacontent_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `order`
--
ALTER TABLE `order`
  MODIFY `order_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `payment`
--
ALTER TABLE `payment`
  MODIFY `payment_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `product`
--
ALTER TABLE `product`
  MODIFY `product_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `product_category`
--
ALTER TABLE `product_category`
  MODIFY `product_category_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `product_detail`
--
ALTER TABLE `product_detail`
  MODIFY `product_detail_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `product_detail_group`
--
ALTER TABLE `product_detail_group`
  MODIFY `product_detail_group_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `rating`
--
ALTER TABLE `rating`
  MODIFY `rating_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `review`
--
ALTER TABLE `review`
  MODIFY `review_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `user`
--
ALTER TABLE `user`
  MODIFY `user_id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `order`
--
ALTER TABLE `order`
  ADD CONSTRAINT `fk_order_invoice1` FOREIGN KEY (`invoice_id`) REFERENCES `invoice` (`invoice_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
