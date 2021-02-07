import { AlphabeticalUsersPipe } from './alphabetical-users.pipe';

describe('AlphabeticalUsersPipe', () => {
  it('create an instance', () => {
    const pipe = new AlphabeticalUsersPipe();
    expect(pipe).toBeTruthy();
  });
});
