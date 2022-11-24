/*
  Warnings:

  - A unique constraint covering the columns `[user_name]` on the table `User` will be added. If there are existing duplicate values, this will fail.
  - Made the column `user_id` on table `bodyinfo` required. This step will fail if there are existing NULL values in that column.

*/
-- DropForeignKey
ALTER TABLE `bodyinfo` DROP FOREIGN KEY `BodyInfo_user_id_fkey`;

-- AlterTable
ALTER TABLE `bodyinfo` MODIFY `user_id` INTEGER NOT NULL;

-- CreateTable
CREATE TABLE `ProgramRec` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `program_id` INTEGER NOT NULL,
    `start` DATETIME(3) NOT NULL,
    `end` DATETIME(3) NOT NULL,
    `sq` DOUBLE NOT NULL,
    `bp` DOUBLE NOT NULL,
    `dl` DOUBLE NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `RoutineRec` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `programrec_id` INTEGER NOT NULL,
    `routine_id` INTEGER NOT NULL,
    `done_reps` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateIndex
CREATE UNIQUE INDEX `User_user_name_key` ON `User`(`user_name`);

-- AddForeignKey
ALTER TABLE `BodyInfo` ADD CONSTRAINT `BodyInfo_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `User`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ProgramRec` ADD CONSTRAINT `ProgramRec_program_id_fkey` FOREIGN KEY (`program_id`) REFERENCES `Program`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `RoutineRec` ADD CONSTRAINT `RoutineRec_programrec_id_fkey` FOREIGN KEY (`programrec_id`) REFERENCES `ProgramRec`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `RoutineRec` ADD CONSTRAINT `RoutineRec_routine_id_fkey` FOREIGN KEY (`routine_id`) REFERENCES `Routine`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
