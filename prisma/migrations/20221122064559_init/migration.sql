/*
  Warnings:

  - You are about to drop the column `bp` on the `programrec` table. All the data in the column will be lost.
  - You are about to drop the column `dl` on the `programrec` table. All the data in the column will be lost.
  - You are about to drop the column `sq` on the `programrec` table. All the data in the column will be lost.
  - Added the required column `benchpress` to the `ProgramRec` table without a default value. This is not possible if the table is not empty.
  - Added the required column `deadlift` to the `ProgramRec` table without a default value. This is not possible if the table is not empty.
  - Added the required column `squat` to the `ProgramRec` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE `programrec` DROP COLUMN `bp`,
    DROP COLUMN `dl`,
    DROP COLUMN `sq`,
    ADD COLUMN `benchpress` DOUBLE NOT NULL,
    ADD COLUMN `deadlift` DOUBLE NOT NULL,
    ADD COLUMN `squat` DOUBLE NOT NULL;
