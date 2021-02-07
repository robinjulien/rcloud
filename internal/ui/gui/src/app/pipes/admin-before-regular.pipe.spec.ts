import { AdminBeforeRegularPipe } from './admin-before-regular.pipe';

describe('AdminBeforeRegularPipe', () => {
  it('create an instance', () => {
    const pipe = new AdminBeforeRegularPipe();
    expect(pipe).toBeTruthy();
  });
});
