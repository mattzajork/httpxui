import { DateTime } from "luxon";

export default class CommonHelper {

    static PB_SERVER = 'http://127.0.0.1:8081';

    /**
     * 
     * @param {Date} dt1 
     * @param {Date} dt2 
     * @returns {string}
     */
    static getDurationFromDates(dt1, dt2) {
        try {
            let dt1DateTime = DateTime.fromFormat(dt1, "yyyy-MM-dd HH:mm:ss.SSS'Z'");
            let dt2DateTime = DateTime.fromFormat(dt2, "yyyy-MM-dd HH:mm:ss.SSS'Z'");
            let diff = dt1DateTime.diff(dt2DateTime);
    
            // calculate seconds
            let totalSeconds = Math.floor(Math.abs(diff) / 1000);
    
            // Calculate hours, minutes, and seconds
            let hours = Math.floor(totalSeconds / 3600);
            let minutes = Math.floor((totalSeconds % 3600) / 60);
            let seconds = totalSeconds % 60;
    
            // Format the time string
            let formattedTime = 
              (hours < 10 ? '0' + hours : hours) + ':' +
              (minutes < 10 ? '0' + minutes : minutes) + ':' +
              (seconds < 10 ? '0' + seconds : seconds);
            return formattedTime;
          } catch(error) {
            return ''
          }
    }

    /** 
     * Adds or replace an object array element by comparing its key value.
     *
     * @param {Array}  objectsArr
     * @param {Object} item
     * @param {Mixed}  [key]
     */
    static pushOrReplaceByKey(objectsArr, item, key = "id") {
        for (let i = objectsArr.length - 1; i >= 0; i--) {
            if (objectsArr[i][key] == item[key]) {
                objectsArr[i] = item; // replace
                return;
            }
        }

        objectsArr.push(item);
    }

    /**
     * Returns single element from objects array by matching its key value.
     *
     * @param  {Array} objectsArr
     * @param  {Mixed} key
     * @param  {Mixed} value
     * @return {Object}
     */
    static findByKey(objectsArr, key, value) {
        objectsArr = Array.isArray(objectsArr) ? objectsArr : [];

        for (let i in objectsArr) {
            if (objectsArr[i][key] == value) {
                return objectsArr[i];
            }
        }

        return null;
    }


    /**
     * Removes single element from objects array by matching an item"s property value.
     *
     * @param {Array}  objectsArr
     * @param {String} key
     * @param {Mixed}  value
     */
    static removeByKey(objectsArr, key, value) {
        for (let i in objectsArr) {
            if (objectsArr[i][key] == value) {
                objectsArr.splice(i, 1);
                break;
            }
        }
    }

    /**
     * Returns a DateTime instance from a date object/string.
     *
     * @param  {String|Date} date
     * @return {DateTime}
     */
    static getDateTime(date) {
        if (typeof date === "string") {
            const formats = {
                19: "yyyy-MM-dd HH:mm:ss",
                23: "yyyy-MM-dd HH:mm:ss.SSS",
                20: "yyyy-MM-dd HH:mm:ss'Z'",
                24: "yyyy-MM-dd HH:mm:ss.SSS'Z'",
            }
            const format = formats[date.length] || formats[19];
            return DateTime.fromFormat(date, format, { zone: "UTC" });
        }

        return DateTime.fromJSDate(date);
    }

    /**
     * Returns formatted datetime string in the UTC timezone.
     *
     * @param  {String|Date} date
     * @param  {String}      [format] The result format (see https://moment.github.io/luxon/#/parsing?id=table-of-tokens)
     * @return {String}
     */
    static formatToUTCDate(date, format = "yyyy-MM-dd HH:mm:ss") {
        return CommonHelper.getDateTime(date).toUTC().toFormat(format);
    }

    /**
     * Returns formatted datetime string in the local timezone.
     *
     * @param  {String|Date} date
     * @param  {String}      [format] The result format (see https://moment.github.io/luxon/#/parsing?id=table-of-tokens)
     * @return {String}
     */
    static formatToLocalDate(date, format = "yyyy-MM-dd HH:mm:ss") {
        return CommonHelper.getDateTime(date).toLocal().toFormat(format);
    }
}