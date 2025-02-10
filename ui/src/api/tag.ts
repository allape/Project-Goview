import Crudy from "@allape/gocrud-react";
import { SERVER_URL } from "@allape/gocrud-react/src/config";
import ITag from "../model/tag.ts";

export const TagCrudy = new Crudy<ITag>(`${SERVER_URL}/tag`);
