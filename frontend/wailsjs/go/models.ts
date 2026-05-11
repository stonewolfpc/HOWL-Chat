export namespace lorebook {
	
	export class Entry {
	    id: string;
	    title: string;
	    content: string;
	    tags: string[];
	    scope: string;
	    type: string;
	    ownerId?: string;
	    enabled: boolean;
	    triggerPhrases: string[];
	    secondaryTriggers: string[];
	    triggerMode: string;
	    triggerDirection: string;
	    triggerFrequency: string;
	    priorityLevel: number;
	    conflictRule: string;
	    visibleToUser: boolean;
	    visibleToOthers: boolean;
	    visibility: string;
	    revealCondition: string;
	    injectionStyle: string;
	    injectionPosition: string;
	    maxLength: string;
	    contextBudget: number;
	    scanWindow: number;
	    recursion: string;
	    linkedCharacters: string[];
	    linkedLoreEntries: string[];
	    inheritance: string;
	    automated: boolean;
	    memoryType: string;
	    triggerConfidence: string;
	    decayRate: string;
	    createdAt: string;
	    // Go type: time
	    lastTriggeredAt?: any;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Entry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.tags = source["tags"];
	        this.scope = source["scope"];
	        this.type = source["type"];
	        this.ownerId = source["ownerId"];
	        this.enabled = source["enabled"];
	        this.triggerPhrases = source["triggerPhrases"];
	        this.secondaryTriggers = source["secondaryTriggers"];
	        this.triggerMode = source["triggerMode"];
	        this.triggerDirection = source["triggerDirection"];
	        this.triggerFrequency = source["triggerFrequency"];
	        this.priorityLevel = source["priorityLevel"];
	        this.conflictRule = source["conflictRule"];
	        this.visibleToUser = source["visibleToUser"];
	        this.visibleToOthers = source["visibleToOthers"];
	        this.visibility = source["visibility"];
	        this.revealCondition = source["revealCondition"];
	        this.injectionStyle = source["injectionStyle"];
	        this.injectionPosition = source["injectionPosition"];
	        this.maxLength = source["maxLength"];
	        this.contextBudget = source["contextBudget"];
	        this.scanWindow = source["scanWindow"];
	        this.recursion = source["recursion"];
	        this.linkedCharacters = source["linkedCharacters"];
	        this.linkedLoreEntries = source["linkedLoreEntries"];
	        this.inheritance = source["inheritance"];
	        this.automated = source["automated"];
	        this.memoryType = source["memoryType"];
	        this.triggerConfidence = source["triggerConfidence"];
	        this.decayRate = source["decayRate"];
	        this.createdAt = source["createdAt"];
	        this.lastTriggeredAt = this.convertValues(source["lastTriggeredAt"], null);
	        this.updatedAt = source["updatedAt"];
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

export namespace types {
	
	export class ModelProfile {
	    Path: string;
	    Family: string;
	    Variant: string;
	    MaxContext: number;
	    UsableContext: number;
	    RopeMode: string;
	    RopeFactor: number;
	    RopeBase: number;
	    Template: string;
	    Tokenizer: string;
	    StopSequences: string[];
	
	    static createFrom(source: any = {}) {
	        return new ModelProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Family = source["Family"];
	        this.Variant = source["Variant"];
	        this.MaxContext = source["MaxContext"];
	        this.UsableContext = source["UsableContext"];
	        this.RopeMode = source["RopeMode"];
	        this.RopeFactor = source["RopeFactor"];
	        this.RopeBase = source["RopeBase"];
	        this.Template = source["Template"];
	        this.Tokenizer = source["Tokenizer"];
	        this.StopSequences = source["StopSequences"];
	    }
	}
	export class PromptSettings {
	    PromptTemplate: string;
	    CustomJinjaTemplate: string;
	    SystemPromptOverride: string;
	    UserPrefix: string;
	    AssistantPrefix: string;
	    StopSequences: string[];
	
	    static createFrom(source: any = {}) {
	        return new PromptSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PromptTemplate = source["PromptTemplate"];
	        this.CustomJinjaTemplate = source["CustomJinjaTemplate"];
	        this.SystemPromptOverride = source["SystemPromptOverride"];
	        this.UserPrefix = source["UserPrefix"];
	        this.AssistantPrefix = source["AssistantPrefix"];
	        this.StopSequences = source["StopSequences"];
	    }
	}
	export class RuntimeSettings {
	    Threads: number;
	    BatchSize: number;
	    ContextSize: number;
	    GPULayers: number;
	    RopeMode: string;
	    RopeFactor: number;
	    RopeBase: number;
	    FlashAttention: boolean;
	    TensorSplit: number[];
	    MainGPU: number;
	    OffloadKQV: boolean;
	    UseMMap: boolean;
	    UseMLock: boolean;
	    VocabOnly: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RuntimeSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Threads = source["Threads"];
	        this.BatchSize = source["BatchSize"];
	        this.ContextSize = source["ContextSize"];
	        this.GPULayers = source["GPULayers"];
	        this.RopeMode = source["RopeMode"];
	        this.RopeFactor = source["RopeFactor"];
	        this.RopeBase = source["RopeBase"];
	        this.FlashAttention = source["FlashAttention"];
	        this.TensorSplit = source["TensorSplit"];
	        this.MainGPU = source["MainGPU"];
	        this.OffloadKQV = source["OffloadKQV"];
	        this.UseMMap = source["UseMMap"];
	        this.UseMLock = source["UseMLock"];
	        this.VocabOnly = source["VocabOnly"];
	    }
	}

}

