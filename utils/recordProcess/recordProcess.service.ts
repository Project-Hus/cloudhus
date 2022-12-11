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
      const al = recordFixed.arm_length;
      const ll = recordFixed.leg_length;
      let al_ = 0;
      let ll_ = 0;
      if ( al === 'long') al_ = 0.75;
      else if (al === 'medium') al_ = 0.5;
      else al_ = 0.25;
      if ( ll === 'long') ll_ = 0.75;
      else if (ll === 'medium') ll_ = 0.5;
      else ll_ = 0.25;
      return {
        sex: recordFixed.sex=='male'?true:false, // true male
        age: recordFixed.age,
        height: recordFixed.height,
        arm_length: al_, 
        leg_length: ll_, 
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
