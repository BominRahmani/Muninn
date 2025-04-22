from rdkit import Chem
from rdkit.Chem import Draw

# Define molecules using SMILES notation
ephedrine = Chem.MolFromSmiles('CN[C@H](C1=CC=CC=C1)[C@@H](O)C')  # Ephedrine
iodoephedrine = Chem.MolFromSmiles('CN[C@H](C1=CC=CC=C1)[C@@H](I)C')  # Intermediate
methamphetamine = Chem.MolFromSmiles('CN[C@H](C1=CC=CC=C1)C')  # Methamphetamine

# Define reaction steps
reaction_steps = [
    (ephedrine, "Ephedrine (C₁₀H₁₅NO)"),
    (iodoephedrine, "Iodoephedrine (Intermediate)"),
    (methamphetamine, "Methamphetamine (C₁₀H₁₅N)")
]

try:
    # Attempt to draw molecules with RDKit
    img = Draw.MolsToGridImage(
        [mol for mol, _ in reaction_steps],
        legends=[desc for _, desc in reaction_steps],
        subImgSize=(400, 200),
        useSVG=False
    )
    img.save("reaction.png")
    print("Reaction diagram saved as reaction.png")

except Exception as e:
    # Text-based fallback if RDKit fails
    print("""
    [Error: RDKit not fully configured. Text Representation Below]
    
    Step 1: Ephedrine
           OH
           |
    C₆H₅-C-CH₂-N-CH₃
           |     |
           CH₃   CH₃
    
    Step 2: Iodoephedrine (Intermediate)
           I
           |
    C₆H₅-C-CH₂-N-CH₃
           |     |
           CH₃   CH₃
    
    Step 3: Methamphetamine (Final Product)
           H
           |
    C₆H₅-C-CH₂-N-CH₃
           |     |
           CH₃   CH₃
    
    Reagents: HI + Red P → Reduction
    """)
