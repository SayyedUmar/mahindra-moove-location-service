-- MySQL dump 10.16  Distrib 10.2.9-MariaDB, for osx10.13 (x86_64)
--
-- Host: localhost    Database: moove_development
-- ------------------------------------------------------
-- Server version	10.2.9-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ar_internal_metadata`
--

DROP TABLE IF EXISTS `ar_internal_metadata`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ar_internal_metadata` (
  `key` varchar(255) NOT NULL,
  `value` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bus_trip_routes`
--

DROP TABLE IF EXISTS `bus_trip_routes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bus_trip_routes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `stop_name` text DEFAULT NULL,
  `stop_address` text DEFAULT NULL,
  `stop_latitude` decimal(10,6) DEFAULT NULL,
  `stop_longitude` decimal(10,6) DEFAULT NULL,
  `stop_order` int(11) DEFAULT NULL,
  `bus_trip_id` int(11) DEFAULT NULL,
  `name` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_bus_trip_routes_on_bus_trip_id` (`bus_trip_id`),
  CONSTRAINT `fk_rails_41366cd05d` FOREIGN KEY (`bus_trip_id`) REFERENCES `bus_trips` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bus_trips`
--

DROP TABLE IF EXISTS `bus_trips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bus_trips` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `status` varchar(255) DEFAULT NULL,
  `route_name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `business_associates`
--

DROP TABLE IF EXISTS `business_associates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `business_associates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_f_name` varchar(255) DEFAULT NULL,
  `admin_m_name` varchar(255) DEFAULT NULL,
  `admin_l_name` varchar(255) DEFAULT NULL,
  `admin_email` varchar(255) DEFAULT NULL,
  `admin_phone` varchar(255) DEFAULT NULL,
  `legal_name` varchar(255) DEFAULT NULL,
  `pan` varchar(255) DEFAULT NULL,
  `tan` varchar(255) DEFAULT NULL,
  `business_type` varchar(255) DEFAULT NULL,
  `service_tax_no` varchar(255) DEFAULT NULL,
  `hq_address` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `standard_price` decimal(10,0) DEFAULT 0,
  `pay_period` int(11) DEFAULT 0,
  `time_on_duty_limit` int(11) DEFAULT 0,
  `distance_limit` int(11) DEFAULT 0,
  `rate_by_time` decimal(10,0) DEFAULT 0,
  `rate_by_distance` decimal(10,0) DEFAULT 0,
  `invoice_frequency` int(11) DEFAULT 0,
  `service_tax_percent` decimal(5,4) DEFAULT 0.0000,
  `swachh_bharat_cess` decimal(5,4) DEFAULT 0.0020,
  `krishi_kalyan_cess` decimal(5,4) DEFAULT 0.0020,
  `logistics_company_id` int(11) DEFAULT NULL,
  `profit_centre` varchar(255) DEFAULT NULL,
  `agreement_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_business_associates_on_logistics_company_id` (`logistics_company_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `configurators`
--

DROP TABLE IF EXISTS `configurators`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `configurators` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `request_type` varchar(255) DEFAULT NULL,
  `value` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `driver_first_pickups`
--

DROP TABLE IF EXISTS `driver_first_pickups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `driver_first_pickups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `trip_id` int(11) DEFAULT NULL,
  `driver_id` int(11) DEFAULT NULL,
  `pickup_time` int(11) DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_driver_first_pickups_on_trip_id` (`trip_id`),
  KEY `index_driver_first_pickups_on_driver_id` (`driver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `driver_requests`
--

DROP TABLE IF EXISTS `driver_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `driver_requests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `request_type` int(11) DEFAULT NULL,
  `reason` int(11) DEFAULT NULL,
  `trip_type` int(11) DEFAULT NULL,
  `request_state` varchar(255) DEFAULT NULL,
  `request_date` datetime DEFAULT NULL,
  `start_date` datetime DEFAULT NULL,
  `end_date` datetime DEFAULT NULL,
  `driver_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `vehicle_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_driver_requests_on_driver_id` (`driver_id`) USING BTREE,
  KEY `index_driver_requests_on_vehicle_id` (`vehicle_id`) USING BTREE,
  CONSTRAINT `fk_rails_6e5b139524` FOREIGN KEY (`driver_id`) REFERENCES `drivers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `drivers`
--

DROP TABLE IF EXISTS `drivers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `drivers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `business_associate_id` int(11) DEFAULT NULL,
  `logistics_company_id` int(11) DEFAULT NULL,
  `site_id` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `badge_number` varchar(255) DEFAULT NULL,
  `badge_issue_date` date DEFAULT NULL,
  `badge_expire_date` date DEFAULT NULL,
  `local_address` varchar(255) DEFAULT NULL,
  `permanent_address` varchar(255) DEFAULT NULL,
  `aadhaar_number` varchar(255) DEFAULT NULL,
  `aadhaar_mobile_number` varchar(255) DEFAULT NULL,
  `licence_number` varchar(255) DEFAULT NULL,
  `licence_validity` date DEFAULT NULL,
  `verified_by_police` tinyint(1) DEFAULT NULL,
  `uniform` tinyint(1) DEFAULT NULL,
  `licence` tinyint(1) DEFAULT NULL,
  `badge` tinyint(1) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `aadhaar_address` varchar(255) DEFAULT NULL,
  `offline_phone` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_drivers_on_business_associate_id` (`business_associate_id`) USING BTREE,
  KEY `index_drivers_on_logistics_company_id` (`logistics_company_id`) USING BTREE,
  KEY `index_drivers_on_site_id` (`site_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=218 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `drivers_shifts`
--

DROP TABLE IF EXISTS `drivers_shifts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `drivers_shifts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `driver_id` int(11) DEFAULT NULL,
  `vehicle_id` int(11) DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `duration` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_drivers_shifts_on_driver_id` (`driver_id`) USING BTREE,
  KEY `index_drivers_shifts_on_vehicle_id` (`vehicle_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5617 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `employee_companies`
--

DROP TABLE IF EXISTS `employee_companies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `employee_companies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `logistics_company_id` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `pan` varchar(255) DEFAULT NULL,
  `tan` varchar(255) DEFAULT NULL,
  `business_type` varchar(255) DEFAULT NULL,
  `service_tax_no` varchar(255) DEFAULT NULL,
  `hq_address` varchar(255) DEFAULT NULL,
  `standard_price` decimal(10,0) DEFAULT 0,
  `pay_period` int(11) DEFAULT 0,
  `time_on_duty_limit` int(11) DEFAULT 0,
  `distance_limit` int(11) DEFAULT 0,
  `rate_by_time` decimal(10,0) DEFAULT 0,
  `rate_by_distance` decimal(10,0) DEFAULT 0,
  `invoice_frequency` int(11) DEFAULT 0,
  `service_tax_percent` decimal(5,4) DEFAULT 0.0000,
  `swachh_bharat_cess` decimal(5,4) DEFAULT 0.0020,
  `krishi_kalyan_cess` decimal(5,4) DEFAULT 0.0020,
  `profit_centre` varchar(255) DEFAULT NULL,
  `agreement_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_employee_companies_on_logistics_company_id` (`logistics_company_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `employee_schedules`
--

DROP TABLE IF EXISTS `employee_schedules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `employee_schedules` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_id` int(11) DEFAULT NULL,
  `day` int(11) DEFAULT NULL,
  `check_in` time DEFAULT NULL,
  `check_out` time DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_employee_schedules_on_employee_id` (`employee_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12013 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `employee_trip_issues`
--

DROP TABLE IF EXISTS `employee_trip_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `employee_trip_issues` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `issue` int(11) DEFAULT NULL,
  `employee_trip_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_employee_trip_issues_on_employee_trip_id` (`employee_trip_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `employee_trips`
--

DROP TABLE IF EXISTS `employee_trips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `employee_trips` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_id` int(11) DEFAULT NULL,
  `trip_id` int(11) DEFAULT NULL,
  `date` datetime DEFAULT NULL,
  `trip_type` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `employee_schedule_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `trip_route_id` int(11) DEFAULT NULL,
  `rating` int(11) DEFAULT NULL,
  `rating_feedback` text DEFAULT NULL,
  `dismissed` tinyint(1) DEFAULT 0,
  `site_id` int(11) DEFAULT NULL,
  `state` int(11) DEFAULT NULL,
  `schedule_date` datetime DEFAULT NULL,
  `zone` int(11) DEFAULT NULL,
  `cluster_error` text DEFAULT NULL,
  `bus_rider` tinyint(1) DEFAULT 0,
  `shift_id` int(11) DEFAULT NULL,
  `is_clustered` tinyint(1) DEFAULT 0,
  `route_order` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_employee_trips_on_employee_id` (`employee_id`) USING BTREE,
  KEY `index_employee_trips_on_trip_id` (`trip_id`) USING BTREE,
  KEY `index_employee_trips_on_trip_route_id` (`trip_route_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=207745 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `employees`
--

DROP TABLE IF EXISTS `employees`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `employees` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_company_id` int(11) DEFAULT NULL,
  `site_id` int(11) DEFAULT NULL,
  `zone_id` int(11) DEFAULT NULL,
  `employee_id` varchar(255) DEFAULT NULL,
  `gender` int(11) DEFAULT NULL,
  `home_address` varchar(255) DEFAULT NULL,
  `home_address_latitude` decimal(10,6) DEFAULT NULL,
  `home_address_longitude` decimal(10,6) DEFAULT NULL,
  `distance_to_site` int(11) DEFAULT NULL,
  `date_of_birth` date DEFAULT NULL,
  `managers_employee_id` varchar(255) DEFAULT NULL,
  `managers_email_id` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `emergency_contact_name` varchar(255) DEFAULT NULL,
  `emergency_contact_phone` varchar(255) DEFAULT NULL,
  `line_manager_id` int(11) DEFAULT NULL,
  `is_guard` tinyint(1) DEFAULT 0,
  `geohash` text DEFAULT NULL,
  `bus_travel` tinyint(1) DEFAULT 0,
  `bus_trip_route_id` int(11) DEFAULT NULL,
  `billing_zone` varchar(255) DEFAULT 'Default',
  PRIMARY KEY (`id`),
  KEY `index_employees_on_employee_company_id` (`employee_company_id`) USING BTREE,
  KEY `index_employees_on_site_id` (`site_id`) USING BTREE,
  KEY `index_employees_on_zone_id` (`zone_id`) USING BTREE,
  KEY `index_employees_on_bus_trip_route_id` (`bus_trip_route_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1717 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `employer_shift_managers`
--

DROP TABLE IF EXISTS `employer_shift_managers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `employer_shift_managers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_company_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `employers`
--

DROP TABLE IF EXISTS `employers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `employers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_company_id` int(11) DEFAULT NULL,
  `legal_name` varchar(255) DEFAULT NULL,
  `pan` varchar(255) DEFAULT NULL,
  `tan` varchar(255) DEFAULT NULL,
  `business_type` varchar(255) DEFAULT NULL,
  `service_tax_no` varchar(255) DEFAULT NULL,
  `hq_address` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_employers_on_employee_company_id` (`employee_company_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ingest_jobs`
--

DROP TABLE IF EXISTS `ingest_jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ingest_jobs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `start_date` date DEFAULT NULL,
  `file_file_name` varchar(255) DEFAULT NULL,
  `file_content_type` varchar(255) DEFAULT NULL,
  `file_file_size` int(11) DEFAULT NULL,
  `file_updated_at` datetime DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `error_file_file_name` varchar(255) DEFAULT NULL,
  `error_file_content_type` varchar(255) DEFAULT NULL,
  `error_file_file_size` int(11) DEFAULT NULL,
  `error_file_updated_at` datetime DEFAULT NULL,
  `failed_row_count` int(11) DEFAULT 0,
  `processed_row_count` int(11) DEFAULT 0,
  `schedule_updated_count` int(11) DEFAULT 0,
  `employee_provisioned_count` int(11) DEFAULT 0,
  `schedule_provisioned_count` int(11) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `index_ingest_jobs_on_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ingest_manifest_jobs`
--

DROP TABLE IF EXISTS `ingest_manifest_jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ingest_manifest_jobs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `status` varchar(255) DEFAULT NULL,
  `failed_row_count` int(11) DEFAULT 0,
  `processed_row_count` int(11) DEFAULT 0,
  `schedule_updated_count` int(11) DEFAULT 0,
  `employee_provisioned_count` int(11) DEFAULT 0,
  `schedule_provisioned_count` int(11) DEFAULT 0,
  `user_id` int(11) DEFAULT NULL,
  `file_file_name` varchar(255) DEFAULT NULL,
  `file_content_type` varchar(255) DEFAULT NULL,
  `file_file_size` int(11) DEFAULT NULL,
  `file_updated_at` datetime DEFAULT NULL,
  `error_file_file_name` varchar(255) DEFAULT NULL,
  `error_file_content_type` varchar(255) DEFAULT NULL,
  `error_file_file_size` int(11) DEFAULT NULL,
  `error_file_updated_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_ingest_manifest_jobs_on_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `invoice_attachments`
--

DROP TABLE IF EXISTS `invoice_attachments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoice_attachments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `invoice_id` int(11) DEFAULT NULL,
  `file_file_name` varchar(255) DEFAULT NULL,
  `file_content_type` varchar(255) DEFAULT NULL,
  `file_file_size` int(11) DEFAULT NULL,
  `file_updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_invoice_attachments_on_invoice_id` (`invoice_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=65 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `invoices`
--

DROP TABLE IF EXISTS `invoices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoices` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `company_type` varchar(255) DEFAULT NULL,
  `company_id` int(11) DEFAULT NULL,
  `date` datetime DEFAULT NULL,
  `start_date` datetime DEFAULT NULL,
  `end_date` datetime DEFAULT NULL,
  `trips_count` int(11) DEFAULT NULL,
  `amount` decimal(12,2) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `status` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_invoices_on_company_type_and_company_id` (`company_type`,`company_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=65 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `line_managers`
--

DROP TABLE IF EXISTS `line_managers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `line_managers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_company_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_line_managers_on_employee_company_id` (`employee_company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `logistics_companies`
--

DROP TABLE IF EXISTS `logistics_companies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `logistics_companies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `pan` varchar(255) DEFAULT NULL,
  `tan` varchar(255) DEFAULT NULL,
  `business_type` varchar(255) DEFAULT NULL,
  `service_tax_no` varchar(255) DEFAULT NULL,
  `hq_address` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `notifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `driver_id` int(11) DEFAULT NULL,
  `employee_id` int(11) DEFAULT NULL,
  `trip_id` int(11) DEFAULT NULL,
  `message` varchar(255) DEFAULT NULL,
  `receiver` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `resolved_status` tinyint(1) DEFAULT 1,
  `call_sid` text DEFAULT NULL,
  `new_notification` tinyint(1) DEFAULT 0,
  `sequence` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_notifications_on_driver_id` (`driver_id`) USING BTREE,
  KEY `index_notifications_on_employee_id` (`employee_id`) USING BTREE,
  KEY `index_notifications_on_trip_id` (`trip_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=171562 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `operator_shift_managers`
--

DROP TABLE IF EXISTS `operator_shift_managers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `operator_shift_managers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `logistics_company_id` int(11) DEFAULT NULL,
  `legal_name` varchar(255) DEFAULT NULL,
  `pan` varchar(255) DEFAULT NULL,
  `tan` varchar(255) DEFAULT NULL,
  `business_type` varchar(255) DEFAULT NULL,
  `service_tax_no` varchar(255) DEFAULT NULL,
  `hq_address` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `operators`
--

DROP TABLE IF EXISTS `operators`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `operators` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `logistics_company_id` int(11) DEFAULT NULL,
  `legal_name` varchar(255) DEFAULT NULL,
  `pan` varchar(255) DEFAULT NULL,
  `tan` varchar(255) DEFAULT NULL,
  `business_type` varchar(255) DEFAULT NULL,
  `service_tax_no` varchar(255) DEFAULT NULL,
  `hq_address` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_operators_on_logistics_company_id` (`logistics_company_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `schema_migrations` (
  `version` varchar(255) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `services`
--

DROP TABLE IF EXISTS `services`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `services` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `site_id` int(11) DEFAULT NULL,
  `service_type` varchar(255) DEFAULT NULL,
  `billing_model` varchar(255) DEFAULT NULL,
  `vary_with_vehicle` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `index_services_on_site_id` (`site_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `shift_times`
--

DROP TABLE IF EXISTS `shift_times`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shift_times` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `shift_manager_id` int(11) DEFAULT NULL,
  `site_id` int(11) DEFAULT NULL,
  `shift_type` int(11) DEFAULT NULL,
  `date` datetime DEFAULT NULL,
  `schedule_date` datetime DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `shift_users`
--

DROP TABLE IF EXISTS `shift_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shift_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `shift_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `shifts`
--

DROP TABLE IF EXISTS `shifts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shifts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `start_time` varchar(255) DEFAULT NULL,
  `end_time` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sites`
--

DROP TABLE IF EXISTS `sites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sites` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `latitude` decimal(10,6) DEFAULT NULL,
  `longitude` decimal(10,6) DEFAULT NULL,
  `employee_company_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_sites_on_employee_company_id` (`employee_company_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `transport_desk_managers`
--

DROP TABLE IF EXISTS `transport_desk_managers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transport_desk_managers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_company_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_transport_desk_managers_on_employee_company_id` (`employee_company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trip_change_requests`
--

DROP TABLE IF EXISTS `trip_change_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trip_change_requests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `request_type` int(11) DEFAULT NULL,
  `reason` int(11) DEFAULT NULL,
  `trip_type` int(11) DEFAULT NULL,
  `request_state` varchar(255) DEFAULT NULL,
  `new_date` datetime DEFAULT NULL,
  `employee_id` int(11) DEFAULT NULL,
  `employee_trip_id` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `shift` tinyint(1) DEFAULT 0,
  `bus_rider` tinyint(1) DEFAULT 0,
  `schedule_date` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_trip_change_requests_on_employee_id` (`employee_id`) USING BTREE,
  KEY `index_trip_change_requests_on_employee_trip_id` (`employee_trip_id`) USING BTREE,
  CONSTRAINT `fk_rails_332c0642cb` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=165 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trip_invoices`
--

DROP TABLE IF EXISTS `trip_invoices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trip_invoices` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `trip_id` int(11) DEFAULT NULL,
  `invoice_id` int(11) DEFAULT NULL,
  `trip_amount` decimal(10,0) DEFAULT 0,
  `trip_penalty` decimal(10,0) DEFAULT 0,
  `trip_toll` decimal(10,0) DEFAULT 0,
  `vehicle_rate_id` int(11) DEFAULT NULL,
  `zone_rate_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_trip_invoices_on_trip_id` (`trip_id`),
  KEY `index_trip_invoices_on_invoice_id` (`invoice_id`),
  KEY `index_trip_invoices_on_vehicle_rate_id` (`vehicle_rate_id`),
  KEY `index_trip_invoices_on_zone_rate_id` (`zone_rate_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trip_locations`
--

DROP TABLE IF EXISTS `trip_locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trip_locations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `trip_id` int(11) DEFAULT NULL,
  `location` text DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `distance` int(11) DEFAULT NULL,
  `speed` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_trip_locations_on_trip_id` (`trip_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10503892 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trip_route_exceptions`
--

DROP TABLE IF EXISTS `trip_route_exceptions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trip_route_exceptions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `trip_route_id` int(11) DEFAULT NULL,
  `date` datetime DEFAULT NULL,
  `exception_type` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `resolved_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5231 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trip_routes`
--

DROP TABLE IF EXISTS `trip_routes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trip_routes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `planned_duration` int(11) DEFAULT NULL,
  `planned_distance` int(11) DEFAULT NULL,
  `planned_route_order` int(11) DEFAULT NULL,
  `planned_start_location` text DEFAULT NULL,
  `planned_end_location` text DEFAULT NULL,
  `employee_trip_id` int(11) DEFAULT NULL,
  `trip_id` int(11) DEFAULT NULL,
  `driver_arrived_date` datetime DEFAULT NULL,
  `on_board_date` datetime DEFAULT NULL,
  `completed_date` datetime DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `scheduled_distance` int(11) DEFAULT NULL,
  `scheduled_duration` int(11) DEFAULT NULL,
  `scheduled_route_order` int(11) DEFAULT NULL,
  `scheduled_start_location` text DEFAULT NULL,
  `scheduled_end_location` text DEFAULT NULL,
  `driver_arrived_location` text DEFAULT NULL,
  `check_in_location` text DEFAULT NULL,
  `drop_off_location` text DEFAULT NULL,
  `missed_location` text DEFAULT NULL,
  `cancel_exception` tinyint(1) DEFAULT 0,
  `cab_type` text DEFAULT NULL,
  `cab_fare` int(11) DEFAULT NULL,
  `cab_driver_name` text DEFAULT NULL,
  `cab_licence_number` text DEFAULT NULL,
  `cab_start_location` text DEFAULT NULL,
  `cab_end_location` text DEFAULT NULL,
  `bus_rider` tinyint(1) DEFAULT 0,
  `bus_stop_name` text DEFAULT NULL,
  `bus_stop_address` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_trip_routes_on_employee_trip_id` (`employee_trip_id`) USING BTREE,
  KEY `index_trip_routes_on_trip_id` (`trip_id`) USING BTREE,
  CONSTRAINT `fk_rails_2069874133` FOREIGN KEY (`trip_id`) REFERENCES `trips` (`id`),
  CONSTRAINT `fk_rails_aa5226250e` FOREIGN KEY (`employee_trip_id`) REFERENCES `employee_trips` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=90631 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trips`
--

DROP TABLE IF EXISTS `trips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `trips` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `driver_id` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `planned_date` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `trip_type` int(11) DEFAULT NULL,
  `planned_approximate_duration` int(11) DEFAULT NULL,
  `start_date` datetime DEFAULT NULL,
  `assign_request_expired_date` datetime DEFAULT NULL,
  `planned_approximate_distance` int(11) DEFAULT NULL,
  `vehicle_id` int(11) DEFAULT NULL,
  `site_id` int(11) DEFAULT NULL,
  `real_duration` int(11) DEFAULT NULL,
  `completed_date` datetime DEFAULT NULL,
  `trip_accept_time` datetime DEFAULT NULL,
  `start_location` text DEFAULT NULL,
  `scheduled_approximate_duration` int(11) DEFAULT NULL,
  `scheduled_approximate_distance` int(11) DEFAULT NULL,
  `scheduled_date` datetime DEFAULT NULL,
  `cancel_status` text DEFAULT NULL,
  `book_ola` tinyint(1) DEFAULT 0,
  `ola_fare` text DEFAULT NULL,
  `bus_rider` tinyint(1) DEFAULT 0,
  `toll` decimal(10,0) DEFAULT 0,
  `penalty` decimal(10,0) DEFAULT 0,
  `amount` decimal(10,0) DEFAULT 0,
  `paid` tinyint(1) DEFAULT 0,
  `is_manual` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `index_trips_on_driver_id` (`driver_id`) USING BTREE,
  KEY `index_trips_on_site_id` (`site_id`) USING BTREE,
  KEY `index_trips_on_vehicle_id` (`vehicle_id`) USING BTREE,
  KEY `index_trips_on_status` (`status`),
  KEY `index_trips_on_scheduled_date` (`scheduled_date`)
) ENGINE=InnoDB AUTO_INCREMENT=27874 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL DEFAULT '',
  `username` varchar(255) DEFAULT NULL,
  `f_name` varchar(255) DEFAULT NULL,
  `m_name` varchar(255) DEFAULT NULL,
  `l_name` varchar(255) DEFAULT NULL,
  `role` int(11) DEFAULT 0,
  `entity_type` varchar(255) DEFAULT NULL,
  `entity_id` int(11) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `encrypted_password` varchar(255) NOT NULL DEFAULT '',
  `reset_password_token` varchar(255) DEFAULT NULL,
  `reset_password_sent_at` datetime DEFAULT NULL,
  `remember_created_at` datetime DEFAULT NULL,
  `sign_in_count` int(11) NOT NULL DEFAULT 0,
  `current_sign_in_at` datetime DEFAULT NULL,
  `last_sign_in_at` datetime DEFAULT NULL,
  `current_sign_in_ip` varchar(255) DEFAULT NULL,
  `last_sign_in_ip` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `tokens` text DEFAULT NULL,
  `provider` varchar(255) NOT NULL DEFAULT 'email',
  `uid` varchar(255) NOT NULL DEFAULT '',
  `avatar_file_name` varchar(255) DEFAULT NULL,
  `avatar_content_type` varchar(255) DEFAULT NULL,
  `avatar_file_size` int(11) DEFAULT NULL,
  `avatar_updated_at` datetime DEFAULT NULL,
  `last_active_time` datetime DEFAULT '2009-01-01 00:00:00',
  `status` int(11) DEFAULT NULL,
  `passcode` varchar(255) DEFAULT NULL,
  `invite_count` int(11) DEFAULT 0,
  `current_location` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_users_on_email` (`email`) USING BTREE,
  UNIQUE KEY `index_users_on_phone` (`phone`) USING BTREE,
  UNIQUE KEY `index_users_on_reset_password_token` (`reset_password_token`) USING BTREE,
  UNIQUE KEY `index_users_on_username` (`username`) USING BTREE,
  KEY `index_users_on_entity_type_and_entity_id` (`entity_type`,`entity_id`) USING BTREE,
  KEY `index_users_on_uid` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=1952 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vehicle_rates`
--

DROP TABLE IF EXISTS `vehicle_rates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vehicle_rates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `service_id` int(11) DEFAULT NULL,
  `vehicle_capacity` int(11) DEFAULT NULL,
  `ac` tinyint(1) DEFAULT 1,
  `cgst` decimal(10,0) DEFAULT 0,
  `sgst` decimal(10,0) DEFAULT 0,
  `overage` tinyint(1) DEFAULT 0,
  `time_on_duty` decimal(10,0) DEFAULT 0,
  `overage_per_hour` decimal(10,0) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `index_vehicle_rates_on_service_id` (`service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vehicles`
--

DROP TABLE IF EXISTS `vehicles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vehicles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `driver_id` int(11) DEFAULT NULL,
  `business_associate_id` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `plate_number` varchar(255) DEFAULT NULL,
  `make` varchar(255) DEFAULT NULL,
  `model` varchar(255) DEFAULT NULL,
  `colour` varchar(255) DEFAULT NULL,
  `driverid` varchar(255) DEFAULT NULL,
  `driver_name` varchar(255) DEFAULT NULL,
  `rc_book_no` varchar(255) DEFAULT NULL,
  `registration_date` date DEFAULT NULL,
  `insurance_date` date DEFAULT NULL,
  `permit_type` varchar(255) DEFAULT NULL,
  `permit_validity_date` date DEFAULT NULL,
  `puc_validity_date` date DEFAULT NULL,
  `fc_validity_date` date DEFAULT NULL,
  `ac` tinyint(1) DEFAULT NULL,
  `seats` int(11) DEFAULT 0,
  `fuel_type` varchar(255) DEFAULT NULL,
  `make_year` int(10) unsigned NOT NULL,
  `induction_date` int(10) unsigned DEFAULT NULL,
  `odometer` int(10) unsigned DEFAULT NULL,
  `spare_type` tinyint(1) DEFAULT NULL,
  `first_aid_kit` tinyint(1) DEFAULT NULL,
  `tyre_condition` varchar(255) DEFAULT NULL,
  `fuel_level` varchar(255) DEFAULT NULL,
  `plate_condition` varchar(255) DEFAULT NULL,
  `device_id` int(10) unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `photo_file_name` varchar(255) DEFAULT NULL,
  `photo_content_type` varchar(255) DEFAULT NULL,
  `photo_file_size` int(11) DEFAULT NULL,
  `photo_updated_at` datetime DEFAULT NULL,
  `status` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_vehicles_on_business_associate_id` (`business_associate_id`) USING BTREE,
  KEY `index_vehicles_on_driver_id` (`driver_id`) USING BTREE,
  KEY `index_vehicles_on_plate_number` (`plate_number`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=173 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `zone_rates`
--

DROP TABLE IF EXISTS `zone_rates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zone_rates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rate` decimal(10,0) DEFAULT 0,
  `guard_rate` decimal(10,0) DEFAULT 0,
  `name` text DEFAULT NULL,
  `vehicle_rate_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_zone_rates_on_vehicle_rate_id` (`vehicle_rate_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `zones`
--

DROP TABLE IF EXISTS `zones`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zones` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` int(11) DEFAULT NULL,
  `latitude` decimal(10,6) DEFAULT NULL,
  `longitude` decimal(10,6) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `site_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_zones_on_site_id` (`site_id`)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-12-18 13:39:52
