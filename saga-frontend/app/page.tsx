'use client';

import { useState } from 'react';

export default function Home() {
  const [item, setItem] = useState('');
  const [amount, setAmount] = useState<string>('');
  const [address, setAddress] = useState('');
  const [simulateShippingFailure, setSimulateShippingFailure] = useState(false);

  const [output, setOutput] = useState<string | null>(null);
  const [error, setError] = useState<boolean>(false);

  const handleSubmit = async () => {
    setOutput(null);
    setError(false);

    try {
      const res = await fetch('/api/orchestrate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ item, amount, address, simulateShippingFailure }),
      });

      const data = await res.json();

      if (!res.ok) {
        setError(true);
        setOutput(data.error || 'Unknown error occurred.');
      } else {
        setError(false);
        setOutput(data.output);
      }
    } catch (err: any) {
      setError(true);
      setOutput(err.message);
    }
  };

  return (
    <main className="max-w-xl mx-auto mt-10">
      <header className="bg-indigo-700 text-white py-4 shadow-md mb-6 rounded">
        <div className="max-w-4xl mx-auto px-4">
          <h1 className="text-2xl font-semibold text-center">
            Tugas Pemrograman Web Lanjut 2 - Tugas 2
          </h1>
          <h2 className="text-lg  text-center">
            Muhammad Zahran Albara - 122140240
          </h2>
        </div>
      </header>
      <input
        type="text"
        value={item}
        onChange={(e) => setItem(e.target.value)}
        placeholder="Item"
        className="w-full p-2 mb-2 border rounded"
      />
      <input
        type="text"
        inputMode="decimal" // memunculkan keyboard angka di HP
        value={amount}
        onChange={(e) => {
          const value = e.target.value;
          // hanya angka dan desimal
          if (/^\d*\.?\d*$/.test(value)) {
            setAmount(value);
          }
        }}
        placeholder="Price"
        className="w-full p-2 mb-2 border rounded"
      />
      <input
        type="text"
        value={address}
        onChange={(e) => setAddress(e.target.value)}
        placeholder="Address"
        className="w-full p-2 mb-2 border rounded"
      />
      <label className="flex items-center space-x-2 mb-4">
        <input
          type="checkbox"
          checked={simulateShippingFailure}
          onChange={(e) => setSimulateShippingFailure(e.target.checked)}
        />
        <span>Simulate Shipping Failure</span>
      </label>
      <button
        onClick={handleSubmit}
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 cursor-pointer"
      >
        Create Order
      </button>

      {output && (
        <pre
          className={`mt-4 p-4 font-mono text-sm rounded whitespace-pre-wrap ${
            error ? 'bg-red-100 text-red-700' : 'bg-black text-green-400'
          }`}
        >
          {output}
        </pre>
      )}
    </main>
  );
}
