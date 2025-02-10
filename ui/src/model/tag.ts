import { IBase } from "@allape/gocrud";

export default interface ITag extends IBase {
  name: string;
  color: string;
}
