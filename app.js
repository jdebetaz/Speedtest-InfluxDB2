require('dotenv').config();

// TODO: Find speedtest-cli module

const schedule = require('node-schedule');
const { InfluxDB, Point } = require('@influxdata/influxdb-client')

function writeData(data) {
    const influx = new InfluxDB({ url: 'http://localhost:808', token })
    const writeApi = influx.getWriteApi('my-org', 'my-bucket')

    const measurement = new Point('speedtest')
        .tag('host', 'server01')
        .floatField('download', data.download)
        .floatField('upload', data.upload)
        .floatField('ping', data.ping)
        .intField('packetLoss', data.packetLoss)
        .timestamp(new Date(data.timestamp));

    writeApi.writePoint(measurement);
    writeApi.close();
}

function worker() {
    (async () => {
        try {
            const result = await SpeedTest({ acceptGdpr: true, acceptLicense: true });
            console.log(result);
        } catch (e) {
            console.error(e);
        }
    })();
}


function bootstrap() {
    const job = schedule.scheduleJob('*/15 * * * *', worker);
}

bootstrap();