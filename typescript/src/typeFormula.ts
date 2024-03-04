export class TypeFormula {
    name: string;
    type_parameters: TypeFormula[];

    constructor(name: string, type_parameters: TypeFormula[]) {
        this.name = name;
        this.type_parameters = type_parameters;
    }

    toString(): string {
        if (this.type_parameters.length > 0) {
            const type_parameters = this.type_parameters.map((type_parameter) => type_parameter.toString()).join(", ");
            return `${this.name}<${type_parameters}>`;
        } else {
            return this.name;
        }
    }
}
