export interface IContainer {
  containerId: number;
  name: string;
  used: number;
  capacity: number;
}

export interface IContainerPreview {
  containerId: number;
  name: string;
}
