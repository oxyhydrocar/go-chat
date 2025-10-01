/**
 * Verifies that setupTests.js correctly wires up Jest with
 * @testing-library/jest-dom matchers.
 *
 * Testing framework: Jest + @testing-library/jest-dom
 */
import './setupTests';

describe('setupTests configuration (Jest + @testing-library/jest-dom)', () => {
  describe('core DOM matchers', () => {
    let element;

    beforeEach(() => {
      element = document.createElement('div');
      element.textContent = 'Core text content';
      document.body.appendChild(element);
    });

    afterEach(() => {
      if (element && element.parentNode) {
        element.parentNode.removeChild(element);
      }
      element = null;
    });

    test('registers toBeInTheDocument matcher', () => {
      expect(typeof expect(element).toBeInTheDocument).toBe('function');
      expect(element).toBeInTheDocument();

      const detached = document.createElement('span');
      expect(detached).not.toBeInTheDocument();
    });

    test('supports text assertions via toHaveTextContent', () => {
      expect(element).toHaveTextContent('Core text content');
      element.textContent = 'Hello World';
      expect(element).toHaveTextContent(/hello/i);
    });

    test('supports class assertions via toHaveClass', () => {
      element.className = 'alpha beta';
      expect(element).toHaveClass('alpha', 'beta');

      expect(() => expect(element).toHaveClass('gamma')).toThrowError(/gamma/);
    });

    test('supports attribute assertions via toHaveAttribute', () => {
      element.setAttribute('data-testid', 'sample');
      expect(element).toHaveAttribute('data-testid', 'sample');
      expect(element).not.toHaveAttribute('aria-label');
    });

    test('supports style and visibility matchers', () => {
      element.style.cssText = 'display: block; opacity: 1;';
      expect(element).toHaveStyle({ display: 'block' });
      expect(element).toBeVisible();

      element.style.display = 'none';
      expect(element).not.toBeVisible();
    });

    test('supports emptiness assertions', () => {
      element.textContent = '';
      expect(element).toBeEmpty();

      element.textContent = 'Restored';
      expect(element).not.toBeEmpty();
    });

    test('handles nullish targets gracefully', () => {
      expect(null).not.toBeInTheDocument();
      expect(undefined).not.toBeInTheDocument();
    });
  });

  describe('form control matchers', () => {
    afterEach(() => {
      document.body.innerHTML = '';
    });

    test('verifies enabled versus disabled states', () => {
      const button = document.createElement('button');
      button.textContent = 'Submit';
      document.body.appendChild(button);

      expect(button).toBeEnabled();
      button.disabled = true;
      expect(button).toBeDisabled();
    });

    test('verifies checked states for inputs', () => {
      const checkbox = document.createElement('input');
      checkbox.type = 'checkbox';
      document.body.appendChild(checkbox);

      expect(checkbox).not.toBeChecked();
      checkbox.checked = true;
      expect(checkbox).toBeChecked();

      const radio = document.createElement('input');
      radio.type = 'radio';
      radio.checked = true;
      document.body.appendChild(radio);
      expect(radio).toBeChecked();
    });

    test('verifies values for textual inputs', () => {
      const input = document.createElement('input');
      input.value = 'hello@example.com';
      document.body.appendChild(input);

      expect(input).toHaveValue('hello@example.com');
      input.value = '';
      expect(input).toHaveValue('');
    });

    test('respects disabled fieldset cascading', () => {
      const fieldset = document.createElement('fieldset');
      fieldset.disabled = true;
      const childInput = document.createElement('input');
      fieldset.appendChild(childInput);
      document.body.appendChild(fieldset);

      expect(childInput).toBeDisabled();
    });
  });

  describe('accessibility-oriented matchers', () => {
    afterEach(() => {
      document.body.innerHTML = '';
    });

    test('exposes accessible names', () => {
      const button = document.createElement('button');
      button.textContent = 'Save';
      document.body.appendChild(button);

      expect(button).toHaveAccessibleName('Save');
    });

    test('exposes accessible descriptions', () => {
      const hint = document.createElement('div');
      hint.id = 'password-hint';
      hint.textContent = 'Must be at least 8 characters.';
      const input = document.createElement('input');
      input.setAttribute('aria-describedby', 'password-hint');

      document.body.appendChild(hint);
      document.body.appendChild(input);

      expect(input).toHaveAccessibleDescription('Must be at least 8 characters.');
    });

    test('reflects validity changes', () => {
      const input = document.createElement('input');
      input.type = 'email';
      input.required = true;
      document.body.appendChild(input);

      expect(input).toBeInvalid();

      input.value = 'user@example.com';
      expect(input).toBeValid();
    });
  });

  describe('matcher diagnostics', () => {
    test('throws informative error messages', () => {
      const element = document.createElement('div');
      element.className = 'existing';
      document.body.appendChild(element);

      expect(() => expect(element).toHaveClass('missing')).toThrowError(/missing/);

      document.body.removeChild(element);
    });
  });

  describe('expect API metadata', () => {
    test('keeps expect.extend available', () => {
      expect(typeof expect.extend).toBe('function');
    });
  });
});