package ar.com.bestprice.buyitnow;

import android.content.Intent;
import android.support.v7.app.AppCompatActivity;
import android.view.View;

import ar.com.bestprice.buyitnow.barcodereader.BarcodeCaptureActivity;

/**
 * Created by ilagiglia on 22/03/2016.
 */
public class AddItemActivity extends AppCompatActivity implements View.OnClickListener {
    private static final int RC_BARCODE_CAPTURE = 9001;

    @Override
    public void onClick(View v) {
       // if (v.getId() == R.id.read_barcode) {
            // launch barcode activity.
            Intent intent = new Intent(this.getApplicationContext(), BarcodeCaptureActivity.class);
            intent.putExtra(BarcodeCaptureActivity.AutoFocus, true);
            intent.putExtra(BarcodeCaptureActivity.UseFlash, false);

            startActivityForResult(intent, RC_BARCODE_CAPTURE);
       // }
    }

}

