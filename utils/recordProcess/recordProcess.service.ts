import { Injectable } from '@nestjs/common';
import { RecordFixed } from 'src/dto/RecordFixed';
import { RecordInput } from 'src/dto/RecordInput';
import { RecordWeekly } from 'src/dto/RecordWeekly';

@Injectable()
export class RecordProcessService {
  /**
   * It attaches the fixed part and variable part of the model input.
   * @param recordFixed Constant values like sex, arm length etc
   * @param recordWeekly Variables like bodyweight etc
   * @returns 
   */
  processRecords(recordFixed:RecordFixed, recordWeekly:RecordWeekly[]): RecordInput[] {
    return recordWeekly.map((r)=>{
      return {
        sex: recordFixed.sex=='male'?true:false, // true male
        age: recordFixed.age,
        height: recordFixed.height,
        arm_length: recordFixed.arm_length, 
        leg_length: recordFixed.leg_length, 
        weight: r.weight,
        fat_rate: r.fat_rate,
        program: r.program, // next program
        squat: r.squat, // weights after previous program
        benchpress: r.benchpress,
        deadlift: r.deadlift,
      }
    });
  }
}
