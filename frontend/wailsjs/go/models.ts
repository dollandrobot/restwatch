export namespace main {
	
	export class Message {
	    publishTime: string;
	    data: string;
	    messageId: string;
	    attributes: Record<string, string>;
	    ExtractedData: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.publishTime = source["publishTime"];
	        this.data = source["data"];
	        this.messageId = source["messageId"];
	        this.attributes = source["attributes"];
	        this.ExtractedData = source["ExtractedData"];
	    }
	}
	export class PubSubMessage {
	    id: string;
	    subscription: string;
	    message: Message;
	    rawMessage: string;
	    // Go type: time
	    receivedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new PubSubMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.subscription = source["subscription"];
	        this.message = this.convertValues(source["message"], Message);
	        this.rawMessage = source["rawMessage"];
	        this.receivedAt = this.convertValues(source["receivedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

