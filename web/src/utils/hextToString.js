export function hexStringToText(hexString) {
    const cleanedHexString = hexString.startsWith('0x') ? hexString.slice(2) : hexString;

    if (cleanedHexString.length % 2 !== 0) {
        throw new Error("Invalid hex string length. Must be even.");
    }

    const bytes = new Uint8Array(cleanedHexString.length / 2);
    for (let i = 0; i < cleanedHexString.length; i += 2) {
        bytes[i / 2] = parseInt(cleanedHexString.substr(i, 2), 16);
    }

    const decoder = new TextDecoder('utf-8');
    const textString = decoder.decode(bytes);

    return textString;
}