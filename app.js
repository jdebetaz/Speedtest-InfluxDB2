require('dotenv').config();

const util = require('node:util');
const exec = util.promisify(require('node:child_process').exec);
const { InfluxDB, Point } = require('@influxdata/influxdb-client')

function writeData(data) {
    const influx = new InfluxDB({ url: process.env.INFLUX_HOST, token: process.env.INFLUX_TOKEN })
    const writeApi = influx.getWriteApi(process.env.INFLUX_ORG, process.env.INFLUX_BUCKET)

    console.log(`Download: ${Math.round(data.download.bytes/1000000)}Mbps - Upload: ${Math.round(data.upload.bytes/1000000)}Mbps - Latency: ${data.ping.latency}ms`)

    const measurement = new Point('speedtest')
        .tag('server', data.server.id)
        .tag('server_name', data.server.name)
        .tag('server_country', data.server.country)
        .floatField('download', data.download.bytes)
        .floatField('upload', data.upload.bytes)
        .floatField('ping', data.ping.latency)
        .stringField('link', data.result.url)
        .timestamp(new Date(data.timestamp));
    writeApi.writePoint(measurement);
    writeApi.close();
}

function worker() {
    (async () => {
        console.log('Running speedtest...');
        try {
            const { stdout, stderr } = await exec('speedtest --accept-license --accept-gdpr --format=json');
            const result = JSON.parse(stdout);
            writeData(result);
        } catch (e) {
            console.error(e);
        }
        console.log(`Next run in ${process.env.APP_INTERVAL} minutes.`)
    })();
}


function bootstrap() {
    setInterval(worker, 1000 * 60 * process.env.APP_INTERVAL);
    worker();
}

bootstrap();