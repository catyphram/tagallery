import { Config } from './models';

declare module '*.json' {
  const value: Config;
  export default value;
}
