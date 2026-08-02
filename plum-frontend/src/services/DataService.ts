import { Data, Dictionary } from '@/dto/Data';

export class DataService {

    getHits(data: Data): number {
        let sum = 0;
        if (data.uris) {
            Object.entries(data.uris).forEach(([uri, uriData]) => {
                if (uriData.statuses) {
                    Object.entries(uriData.statuses).forEach(([status, statusData]) => {
                        sum += statusData.hits;
                    });
                }
            });
        }
        return sum;
    }

    getVisits(data: Data): number {
        return data.visits;
    }

    getBytesSent(data: Data): number {
        let sum = 0;
        if (data.uris) {
            Object.entries(data.uris).forEach(([uri, uriData]) => {
                sum += uriData.bytes;
            });
        }
        return sum;
    }

    getStatusMapping(data: Data): Dictionary<Dictionary<number>> {
        if (data.uris) {
            return Object.entries(data.uris).reduce((acc, [uri, uriData]) => {
                Object.entries(uriData.statuses).forEach(([status, statusData]) => {
                    if (!(status in acc)) {
                        acc[status] = {};
                    }
                    acc[status][uri] = (acc[status][uri] || 0) + statusData.hits;
                });
                return acc;
            }, {});
        }
        return {};
    }

}

